package service

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/repository"
	"astigo/pkg/model"
	"context"
)

var (
	_ handler.IFooHandler = (*FooService)(nil)
)

type FooService struct {
	repo repository.IFooRepository
}

func NewService(repo repository.IFooRepository) *FooService {
	return &FooService{repo: repo}
}

func (s *FooService) Get(ctx context.Context, id string) (*model.Foo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FooService) Register(ctx context.Context, input model.Foo) error {
	return s.repo.Create(ctx, &input)
}
