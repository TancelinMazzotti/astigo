package handler

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooHandler = (*MockFooHandler)(nil)
)

type MockFooHandler struct {
	mock.Mock
}

func (m *MockFooHandler) GetAll(ctx context.Context, pagination PaginationInput) ([]model.Foo, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]model.Foo), args.Error(1)
}

func (m *MockFooHandler) GetByID(ctx context.Context, id int) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooHandler) Create(ctx context.Context, input FooCreateInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooHandler) Update(ctx context.Context, input FooUpdateInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooHandler) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
