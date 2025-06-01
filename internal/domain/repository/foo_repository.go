package repository

import (
	"astigo/pkg/dto"
	"context"
)

type IFooRepository interface {
	FindAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error)
	FindByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error)
	Create(ctx context.Context, foo dto.FooRequestCreateDto) error
	Update(ctx context.Context, foo dto.FooRequestUpdateDto) error
	DeleteByID(ctx context.Context, id int) error
}
