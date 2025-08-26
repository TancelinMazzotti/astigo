package cache

import (
	"astigo/internal/domain/model"
	"context"

	"github.com/google/uuid"

	"time"
)

// IFooCache defines a contract for caching operations related to Foo entities.
// GetByID retrieves a Foo entity from the cache by its UUID. Returns an error if the operation fails.
// Set stores a Foo entity in the cache with the specified expiration duration. Returns an error if the operation fails.
// DeleteByID removes a Foo entity from the cache using its UUID. Returns an error if the operation fails.
type IFooCache interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Set(ctx context.Context, foo *model.Foo, expiration time.Duration) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
