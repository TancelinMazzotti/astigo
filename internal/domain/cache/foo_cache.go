package cache

import (
	"astigo/internal/domain/model"
	"context"

	"time"
)

type IFooCache interface {
	GetByID(ctx context.Context, id int) (*model.Foo, error)
	Set(ctx context.Context, foo model.Foo, expiration time.Duration) error
}
