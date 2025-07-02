package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

const FooCacheExpiration = time.Minute * 15

var (
	_ handler.IFooHandler = (*FooService)(nil)
)

// FooService provides business logic around Foo entities, integrating data access, caching, and messaging capabilities.
type FooService struct {
	logger    *zap.Logger
	repo      repository.IFooRepository
	cache     cache.IFooCache
	messaging messaging.IFooMessaging
}

// GetAll retrieves a list of Foo entities based on the provided input criteria and returns an error if retrieval fails.
func (s *FooService) GetAll(ctx context.Context, input handler.FooReadListInput) ([]*model.Foo, error) {
	foos, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("fail to find all foo: %w", err)
	}

	return foos, nil
}

// GetByID retrieves a Foo entity by its ID, using a cache-first approach and falling back to the repository if needed.
func (s *FooService) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	foo, err := s.cache.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn("fail to find foo by id from cache", zap.Error(err))
	}

	if foo == nil {
		foo, err = s.repo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fail to find foo by id: %w", err)
		}

		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			s.logger.Warn("fail to create foo in cache: %w", zap.Error(err))
		}
	}

	return foo, nil
}

// Create creates a new Foo entity, stores it in the repository, and updates related cache and messaging.
func (s *FooService) Create(ctx context.Context, input handler.FooCreateInput) (*model.Foo, error) {
	foo := &model.Foo{
		Id:     uuid.New(),
		Label:  input.Label,
		Secret: input.Secret,
		Value:  input.Value,
		Weight: input.Weight,
	}

	var validate = validator.New()
	if err := validate.Struct(foo); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Create(ctx, foo); err != nil {
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
		return nil, fmt.Errorf("fail to publish foo created: %w", errMessaging)
	}

	return foo, nil
}

// Update applies partial updates to an existing Foo entity based on the provided input and propagates changes across systems.
// It retrieves the entity by ID, merges changes, updates the repository, cache, and publishes an event.
func (s *FooService) Update(ctx context.Context, input handler.FooUpdateInput) error {
	foo, err := s.repo.FindByID(ctx, input.Id)
	if err != nil {
		return fmt.Errorf("fail to get foo by id: %w", err)
	}

	if err := input.Merge(foo); err != nil {
		return fmt.Errorf("fail to merge input: %w", err)
	}

	var validate = validator.New()
	if err := validate.Struct(foo); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Update(ctx, foo); err != nil {
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
		return fmt.Errorf("fail to publish foo updated: %w", errMessaging)
	}

	return nil
}

// DeleteByID removes a Foo entity by its ID, updates the cache, and publishes a deletion event. Returns an error if any step fails.
func (s *FooService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
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
