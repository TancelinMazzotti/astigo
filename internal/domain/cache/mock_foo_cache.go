package cache

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

var (
	_ IFooCache = (*MockFooCache)(nil)
)

type MockFooCache struct {
	mock.Mock
}

func (m *MockFooCache) GetByID(ctx context.Context, id int) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooCache) Set(ctx context.Context, foo model.Foo, expiration time.Duration) error {
	args := m.Called(ctx, foo, expiration)
	return args.Error(0)
}
