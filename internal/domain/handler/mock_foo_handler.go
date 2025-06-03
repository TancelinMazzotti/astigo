package handler

import (
	"astigo/pkg/dto"
	"context"
	"github.com/stretchr/testify/mock"
)

var (
	_ IFooHandler = (*MockFooHandler)(nil)
)

type MockFooHandler struct {
	mock.Mock
}

func (m *MockFooHandler) GetAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]dto.FooResponseReadDto), args.Error(1)
}

func (m *MockFooHandler) GetByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.FooResponseReadDto), args.Error(1)
}

func (m *MockFooHandler) Create(ctx context.Context, input dto.FooRequestCreateDto) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooHandler) Update(ctx context.Context, input dto.FooRequestUpdateDto) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockFooHandler) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
