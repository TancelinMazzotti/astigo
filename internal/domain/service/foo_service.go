package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"fmt"
)

var (
	_ handler.IFooHandler = (*FooService)(nil)
)

type FooService struct {
	repo      repository.IFooRepository
	cache     cache.IFooCache
	messaging messaging.IFooMessaging
}

func (s *FooService) GetAll(ctx context.Context, pagination handler.PaginationInput) ([]model.Foo, error) {
	foos, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("fail to find all foo: %w", err)
	}

	return foos, nil
}

func (s *FooService) GetByID(ctx context.Context, id int) (*model.Foo, error) {
	foo, err := s.cache.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to find foo by id from cache: %w", err)
	}

	if foo == nil {
		foo, err = s.repo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fail to find foo by id: %w", err)
		}

		if err := s.cache.Set(ctx, *foo, 0); err != nil {
			return nil, fmt.Errorf("fail to create foo in cache: %w", err)
		}
	}

	return foo, nil
}

func (s *FooService) Create(ctx context.Context, input handler.FooCreateInput) error {
	if err := s.repo.Create(ctx, input); err != nil {
		return fmt.Errorf("fail to create foo: %w", err)
	}

	return nil
}

func (s *FooService) Update(ctx context.Context, input handler.FooUpdateInput) error {
	if err := s.repo.Update(ctx, input); err != nil {
		return fmt.Errorf("fail to update foo: %w", err)
	}

	return nil
}

func (s *FooService) DeleteByID(ctx context.Context, id int) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("fail to delete foo by id: %w", err)
	}

	return nil
}

func NewFooService(repo repository.IFooRepository, cache cache.IFooCache, messaging messaging.IFooMessaging) *FooService {
	return &FooService{
		repo:      repo,
		cache:     cache,
		messaging: messaging,
	}
}
