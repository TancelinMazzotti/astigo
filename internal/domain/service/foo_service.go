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
	cache     cache.IFooCahe
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
	foo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to find foo by id: %w", err)
	}
	return foo, nil
}

func (s *FooService) Create(ctx context.Context, input handler.FooCreateInput) error {
	return s.repo.Create(ctx, input)
}

func (s *FooService) Update(ctx context.Context, input handler.FooUpdateInput) error {
	return s.repo.Update(ctx, input)
}

func (s *FooService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

func NewService(repo repository.IFooRepository, cache cache.IFooCahe, messaging messaging.IFooMessaging) *FooService {
	return &FooService{
		repo:      repo,
		cache:     cache,
		messaging: messaging,
	}
}
