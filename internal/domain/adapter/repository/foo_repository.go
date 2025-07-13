package repository

import (
	"astigo/internal/domain/adapter/data"
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooRepository = (*MockFooRepository)(nil)
)

// IFooRepository represents a contract for interacting with Foo data storage.
// FindAll retrieves a paginated list of Foo entities from the repository.
// FindByID fetches a Foo entity by its unique identifier.
// Create adds a new Foo entity to the repository.
// Update modifies an existing Foo entity in the repository.
// DeleteByID removes a Foo entity by its unique identifier from the repository.
type IFooRepository interface {
	FindAll(ctx context.Context, pagination data.FooReadListInput) ([]*model.Foo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Create(ctx context.Context, foo *model.Foo) error
	Update(ctx context.Context, foo *model.Foo) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type MockFooRepository struct {
	mock.Mock
}

func (m *MockFooRepository) FindAll(ctx context.Context, pagination data.FooReadListInput) ([]*model.Foo, error) {
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
