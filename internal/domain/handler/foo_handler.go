package handler

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
)

type IFooHandler interface {
	GetAll(ctx context.Context, input FooReadListInput) ([]model.Foo, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Create(ctx context.Context, input FooCreateInput) (*model.Foo, error)
	Update(ctx context.Context, input FooUpdateInput) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type FooReadListInput struct {
	Offset int
	Limit  int
}

type FooReadInput struct {
	Id uuid.UUID
}

type FooCreateInput struct {
	Label  string
	Secret string
}

type FooUpdateInput struct {
	Id     uuid.UUID
	Label  string
	Secret string
}

func (f FooUpdateInput) Merge(foo *model.Foo) error {
	foo.Label = f.Label
	foo.Secret = f.Secret

	return nil
}

type FooDeleteInput struct {
	Id uuid.UUID
}
