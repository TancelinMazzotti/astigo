package repository

import (
	"astigo/pkg/model"
	"context"
)

type IFooRepository interface {
	FindByID(ctx context.Context, id string) (*model.Foo, error)
	Create(ctx context.Context, foo *model.Foo) error
}
