package handler

import (
	"astigo/pkg/dto"
	"context"
)

type IFooHandler interface {
	GetAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error)
	GetByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error)
	Create(ctx context.Context, input dto.FooRequestCreateDto) error
	Update(ctx context.Context, input dto.FooRequestUpdateDto) error
	DeleteByID(ctx context.Context, id int) error
}
