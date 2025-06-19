package cache

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"

	"time"
)

type IFooCache interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Set(ctx context.Context, foo model.Foo, expiration time.Duration) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
