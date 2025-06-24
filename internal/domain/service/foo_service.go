package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

const FooCacheExpiration = time.Minute * 15

var (
	_ handler.IFooHandler = (*FooService)(nil)
)

type FooService struct {
	logger    *zap.Logger
	repo      repository.IFooRepository
	cache     cache.IFooCache
	messaging messaging.IFooMessaging
}

func (s *FooService) GetAll(ctx context.Context, input handler.FooReadListInput) ([]model.Foo, error) {
	foos, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("fail to find all foo: %w", err)
	}

	return foos, nil
}

func (s *FooService) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	foo, err := s.cache.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to find foo by id from cache: %w", err)
	}

	if foo == nil {
		foo, err = s.repo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fail to find foo by id: %w", err)
		}

		if err := s.cache.Set(ctx, *foo, FooCacheExpiration); err != nil {
			return nil, fmt.Errorf("fail to create foo in cache: %w", err)
		}
	}

	return foo, nil
}

func (s *FooService) Create(ctx context.Context, input handler.FooCreateInput) (*model.Foo, error) {
	foo := model.Foo{
		Id:     uuid.New(),
		Label:  input.Label,
		Secret: input.Secret,
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

	return &foo, nil
}

func (s *FooService) Update(ctx context.Context, input handler.FooUpdateInput) error {
	foo, err := s.repo.FindByID(ctx, input.Id)
	if err != nil {
		return fmt.Errorf("fail to get foo by id: %w", err)
	}

	if err := input.Merge(foo); err != nil {
		return fmt.Errorf("fail to merge input: %w", err)
	}

	if err := s.repo.Update(ctx, *foo); err != nil {
		return fmt.Errorf("fail to update foo: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.Set(ctx, *foo, FooCacheExpiration); err != nil {
			s.logger.Warn("fail to update foo in cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		errMessaging = s.messaging.PublishFooUpdated(ctx, *foo)
	}()

	wg.Wait()

	if errMessaging != nil {
		return fmt.Errorf("fail to publish foo updated: %w", errMessaging)
	}

	return nil
}

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

func NewFooService(logger *zap.Logger, repo repository.IFooRepository, cache cache.IFooCache, messaging messaging.IFooMessaging) *FooService {
	return &FooService{
		logger:    logger,
		repo:      repo,
		cache:     cache,
		messaging: messaging,
	}
}
