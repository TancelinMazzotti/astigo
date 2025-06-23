package messaging

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooMessaging = (*MockFooMessaging)(nil)
)

type MockFooMessaging struct {
	mock.Mock
}

func (m *MockFooMessaging) PublishFooCreated(ctx context.Context, foo model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooMessaging) PublishFooUpdated(ctx context.Context, foo model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooMessaging) PublishFooDeleted(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
