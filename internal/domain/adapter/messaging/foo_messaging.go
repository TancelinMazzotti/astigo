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

// IFooMessaging defines a contract for publishing events related to Foo entities.
// PublishFooCreated sends a message when a Foo entity is created.
// PublishFooUpdated sends a message when a Foo entity is updated.
// PublishFooDeleted sends a message when a Foo entity is deleted.
type IFooMessaging interface {
	PublishFooCreated(ctx context.Context, foo *model.Foo) error
	PublishFooUpdated(ctx context.Context, foo *model.Foo) error
	PublishFooDeleted(ctx context.Context, id uuid.UUID) error
}

type MockFooMessaging struct {
	mock.Mock
}

func (m *MockFooMessaging) PublishFooCreated(ctx context.Context, foo *model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooMessaging) PublishFooUpdated(ctx context.Context, foo *model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooMessaging) PublishFooDeleted(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
