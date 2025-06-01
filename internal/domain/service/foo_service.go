package service

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/repository"
	"astigo/pkg/dto"
	"context"
)

var (
	_ handler.IFooHandler = (*FooService)(nil)
)

type FooService struct {
	repo repository.IFooRepository
}

func (s *FooService) GetAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error) {
	return s.repo.FindAll(ctx, pagination)
}

func (s *FooService) GetByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FooService) Create(ctx context.Context, input dto.FooRequestCreateDto) error {
	return s.repo.Create(ctx, input)
}

func (s *FooService) Update(ctx context.Context, input dto.FooRequestUpdateDto) error {
	return s.repo.Update(ctx, input)
}

func (s *FooService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

func NewService(repo repository.IFooRepository) *FooService {
	return &FooService{repo: repo}
}
