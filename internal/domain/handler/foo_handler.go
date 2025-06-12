package handler

import (
	"astigo/internal/domain/model"
	"context"
)

type IFooHandler interface {
	GetAll(ctx context.Context, pagination PaginationInput) ([]model.Foo, error)
	GetByID(ctx context.Context, id int) (*model.Foo, error)
	Create(ctx context.Context, input FooCreateInput) error
	Update(ctx context.Context, input FooUpdateInput) error
	DeleteByID(ctx context.Context, id int) error
}

type FooReadInput struct {
	Id int
}

type FooReadOutput struct {
	Id    int
	Label string
}

type FooCreateInput struct {
	Label  string
	Secret string
}

type FooUpdateInput struct {
	Id     int
	Label  string
	Secret string
}

type FooDeleteInput struct {
	Id int
}
