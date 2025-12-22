package mockpublisher

import (
	"context"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/publisher"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	_ publisher.IFooPublisher = (*MockFooPublisher)(nil)
)

type MockFooPublisher struct {
	mock.Mock
}

func (m *MockFooPublisher) PublishCreated(ctx context.Context, foo *model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooPublisher) PublishUpdated(ctx context.Context, foo *model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooPublisher) PublishDeleted(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
