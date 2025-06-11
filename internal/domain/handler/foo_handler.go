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
	Id int `uri:"id" binding:"required,numeric"`
}

type FooReadOutput struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type FooCreateInput struct {
	Label  string `json:"label"`
	Secret string `json:"secret"`
}

type FooUpdateInput struct {
	Id     int    `json:"id"`
	Label  string `json:"label"`
	Secret string `json:"secret"`
}

type FooDeleteInput struct {
	Id int `uri:"id" binding:"required,numeric"`
}
