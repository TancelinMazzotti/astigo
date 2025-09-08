package service

import (
	"context"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/data"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	_ service.IFooService = (*MockFooService)(nil)
)

type MockFooService struct {
	mock.Mock
}

func (m *MockFooService) GetAll(ctx context.Context, pagination data.FooReadListInput) ([]*model.Foo, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]*model.Foo), args.Error(1)
}

func (m *MockFooService) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooService) Create(ctx context.Context, input data.FooCreateInput) (*model.Foo, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooService) Update(ctx context.Context, input data.IFooUpdateMerger) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
