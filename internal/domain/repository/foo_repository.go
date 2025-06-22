package repository

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
)

type IFooRepository interface {
	FindAll(ctx context.Context, pagination handler.FooReadListInput) ([]model.Foo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Create(ctx context.Context, foo model.Foo) error
	Update(ctx context.Context, foo model.Foo) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
