package handler

import (
	"astigo/pkg/model"
	"context"
)

type IFooHandler interface {
	Get(ctx context.Context, id string) (*model.Foo, error)
	Register(ctx context.Context, input model.Foo) error
}
