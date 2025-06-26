package repository

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooRepository = (*MockFooRepository)(nil)
)

type MockFooRepository struct {
	mock.Mock
}

func (m *MockFooRepository) FindAll(ctx context.Context, pagination handler.FooReadListInput) ([]*model.Foo, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]*model.Foo), args.Error(1)
}

func (m *MockFooRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Foo), args.Error(1)
}

func (m *MockFooRepository) Create(ctx context.Context, foo *model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooRepository) Update(ctx context.Context, foo *model.Foo) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
