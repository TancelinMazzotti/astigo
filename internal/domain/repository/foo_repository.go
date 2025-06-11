package repository

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
)

type IFooRepository interface {
	FindAll(ctx context.Context, pagination handler.PaginationInput) ([]model.Foo, error)
	FindByID(ctx context.Context, id int) (*model.Foo, error)
	Create(ctx context.Context, foo handler.FooCreateInput) error
	Update(ctx context.Context, foo handler.FooUpdateInput) error
	DeleteByID(ctx context.Context, id int) error
}
