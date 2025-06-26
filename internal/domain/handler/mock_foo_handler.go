package handler

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooHandler = (*MockFooHandler)(nil)
)

type MockFooHandler struct {
	mock.Mock
}

func (m *MockFooHandler) GetAll(ctx context.Context, pagination FooReadListInput) ([]*model.Foo, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]*model.Foo), args.Error(1)
}

func (m *MockFooHandler) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooHandler) Create(ctx context.Context, input FooCreateInput) (*model.Foo, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooHandler) Update(ctx context.Context, input FooUpdateInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooHandler) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
