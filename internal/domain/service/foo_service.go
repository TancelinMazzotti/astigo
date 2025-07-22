package service

import (
	"astigo/internal/domain/adapter/cache"
	"astigo/internal/domain/adapter/data"
	"astigo/internal/domain/adapter/messaging"
	"astigo/internal/domain/adapter/repository"
	"astigo/internal/domain/model"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"sync"
	"time"
)

const FooCacheExpiration = time.Minute * 15

var (
	_ IFooService = (*MockFooService)(nil)
	_ IFooService = (*FooService)(nil)
)

// IFooService defines the interface for handling operations related to Foo entities.
// GetAll retrieves a list of Foo entities based on the provided input.
// GetByID fetches a Foo entity by its unique identifier.
// Create adds a new Foo entity based on the provided input and returns the created instance.
// Update modifies an existing Foo entity based on the provided input.
// DeleteByID removes a Foo entity identified by its unique identifier.
type IFooService interface {
	GetAll(ctx context.Context, input data.FooReadListInput) ([]*model.Foo, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Create(ctx context.Context, input data.FooCreateInput) (*model.Foo, error)
	Update(ctx context.Context, input data.IFooUpdateMerger) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type MockFooService struct {
	mock.Mock
}

func (m *MockFooService) GetAll(ctx context.Context, pagination data.FooReadListInput) ([]*model.Foo, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]*model.Foo), args.Error(1)
}

func (m *MockFooService) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooService) Create(ctx context.Context, input data.FooCreateInput) (*model.Foo, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooService) Update(ctx context.Context, input data.IFooUpdateMerger) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// FooService provides business logic around Foo entities, integrating data access, caching, and messaging capabilities.
type FooService struct {
	logger    *zap.Logger
	repo      repository.IFooRepository
	cache     cache.IFooCache
	messaging messaging.IFooMessaging
}

// GetAll retrieves a list of Foo entities based on the provided input criteria and returns an error if retrieval fails.
func (s *FooService) GetAll(ctx context.Context, input data.FooReadListInput) ([]*model.Foo, error) {
	foos, err := s.repo.FindAll(ctx, input)
	if err != nil {
		s.logger.Debug("fail to find all foo", zap.Error(err))
		return nil, fmt.Errorf("fail to find all foo: %w", err)
	}

	return foos, nil
}

// GetByID retrieves a Foo entity by its ID, using a cache-first approach and falling back to the repository if needed.
func (s *FooService) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	foo, err := s.cache.GetByID(ctx, id)
	if err != nil {
		s.logger.Debug("fail to find foo by id from cache", zap.Error(err))
	}

	if foo == nil {
		foo, err = s.repo.FindByID(ctx, id)
		if err != nil {
			s.logger.Debug("fail to find foo by id", zap.Error(err))
			return nil, fmt.Errorf("fail to find foo by id: %w", err)
		}

		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			s.logger.Warn("fail to create foo in cache: %w", zap.Error(err))
		}
	}

	return foo, nil
}

// Create creates a new Foo entity, stores it in the repository, and updates related cache and messaging.
func (s *FooService) Create(ctx context.Context, input data.FooCreateInput) (*model.Foo, error) {
	foo := &model.Foo{
		Id:     uuid.New(),
		Label:  input.Label,
		Secret: input.Secret,
		Value:  input.Value,
		Weight: input.Weight,
	}

	var validate = validator.New()
	if err := validate.Struct(foo); err != nil {
		s.logger.Debug("invalid input", zap.Error(err))
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Create(ctx, foo); err != nil {
		s.logger.Debug("fail to create foo", zap.Error(err))
		return nil, fmt.Errorf("fail to create foo: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			s.logger.Warn("fail to create foo in cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		errMessaging = s.messaging.PublishFooCreated(ctx, foo)
	}()

	wg.Wait()

	if errMessaging != nil {
		s.logger.Debug("fail to publish foo created", zap.Error(errMessaging))
		return nil, fmt.Errorf("fail to publish foo created: %w", errMessaging)
	}

	return foo, nil
}

// Update applies partial updates to an existing Foo entity based on the provided input and propagates changes across systems.
// It retrieves the entity by ID, merges changes, updates the repository, cache, and publishes an event.
func (s *FooService) Update(ctx context.Context, input data.IFooUpdateMerger) error {
	foo, err := s.repo.FindByID(ctx, input.GetID())
	if err != nil {
		s.logger.Debug("fail to find foo by id", zap.Error(err))
		return fmt.Errorf("fail to get foo by id: %w", err)
	}

	if err := input.Merge(foo); err != nil {
		s.logger.Debug("fail to merge input", zap.Error(err))
		return fmt.Errorf("fail to merge input: %w", err)
	}

	var validate = validator.New()
	if err := validate.Struct(foo); err != nil {
		s.logger.Debug("invalid input", zap.Error(err))
		return fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Update(ctx, foo); err != nil {
		s.logger.Debug("fail to update foo", zap.Error(err))
		return fmt.Errorf("fail to update foo: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			s.logger.Warn("fail to update foo in cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		errMessaging = s.messaging.PublishFooUpdated(ctx, foo)
	}()

	wg.Wait()

	if errMessaging != nil {
		s.logger.Debug("fail to publish foo updated", zap.Error(errMessaging))
		return fmt.Errorf("fail to publish foo updated: %w", errMessaging)
	}

	return nil
}

// DeleteByID removes a Foo entity by its ID, updates the cache, and publishes a deletion event. Returns an error if any step fails.
func (s *FooService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		s.logger.Debug("fail to delete foo by id", zap.Error(err))
		return fmt.Errorf("fail to delete foo by id: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.DeleteByID(ctx, id); err != nil {
			s.logger.Warn("fail to delete foo by id from cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		errMessaging = s.messaging.PublishFooDeleted(ctx, id)
	}()

	wg.Wait()

	if errMessaging != nil {
		s.logger.Debug("fail to publish foo deleted", zap.Error(errMessaging))
		return fmt.Errorf("fail to publish foo deleted: %w", errMessaging)
	}

	return nil
}

// NewFooService initializes a new instance of FooService with the provided logger, repository, cache, and messaging dependencies.
func NewFooService(logger *zap.Logger, repo repository.IFooRepository, cache cache.IFooCache, messaging messaging.IFooMessaging) *FooService {
	return &FooService{
		logger:    logger,
		repo:      repo,
		cache:     cache,
		messaging: messaging,
	}
}
