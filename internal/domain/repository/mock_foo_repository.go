package repository

import (
	"astigo/pkg/dto"
	"context"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooRepository = (*MockFooRepository)(nil)
)

type MockFooRepository struct {
	mock.Mock
}

func (m *MockFooRepository) FindAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]dto.FooResponseReadDto), args.Error(1)
}

func (m *MockFooRepository) FindByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.FooResponseReadDto), args.Error(1)
}

func (m *MockFooRepository) Create(ctx context.Context, foo dto.FooRequestCreateDto) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooRepository) Update(ctx context.Context, foo dto.FooRequestUpdateDto) error {
	args := m.Called(ctx, foo)
	return args.Error(0)
}

func (m *MockFooRepository) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
