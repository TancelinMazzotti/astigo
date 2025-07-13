package cache

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"time"
)

var (
	_ IFooCache = (*MockFooCache)(nil)
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

type MockFooCache struct {
	mock.Mock
}

func (m *MockFooCache) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooCache) Set(ctx context.Context, foo *model.Foo, expiration time.Duration) error {
	args := m.Called(ctx, foo, expiration)
	return args.Error(0)
}

func (m *MockFooCache) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
