package grpc

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"astigo/pkg/proto"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFooService_List(t *testing.T) {
	testCases := []struct {
		name    string
		request *proto.ListFoosRequest

		mockRequest  handler.PaginationInput
		mockResponse []model.Foo
		mockError    error

		expectedError error
		expectedCount int
	}{
		{
			name: "Success Case",
			request: &proto.ListFoosRequest{
				Offset: 0,
				Limit:  10,
			},

			mockRequest: handler.PaginationInput{Offset: 0, Limit: 10},
			mockResponse: []model.Foo{
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "Foo1"},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "Foo2"},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "Foo3"},
			},
			mockError: nil,

			expectedCount: 3,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("GetAll", mock.Anything, testCase.mockRequest).Return(testCase.mockResponse, testCase.mockError)

			service := NewFooService(mockHandler)

			resp, err := service.List(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.expectedError.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Len(t, resp.Foos, testCase.expectedCount)
			}

			mockHandler.AssertExpectations(t)
		})
	}
}

func TestFooService_Get(t *testing.T) {
	testCases := []struct {
		name    string
		request *proto.GetFooRequest

		mockRequest  uuid.UUID
		mockResponse *model.Foo
		mockError    error

		expectedError  error
		expectedResult *proto.FooResponse
	}{
		{
			name: "Success Case",
			request: &proto.GetFooRequest{
				Id: "20000000-0000-0000-0000-000000000001",
			},
			mockRequest: uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockResponse: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "Foo1",
				Secret: "secret1",
			},
			mockError: nil,

			expectedError: nil,
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:    "20000000-0000-0000-0000-000000000001",
					Label: "Foo1",
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("GetByID", mock.Anything, testCase.mockRequest).Return(testCase.mockResponse, testCase.mockError)

			service := NewFooService(mockHandler)

			resp, err := service.Get(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.expectedError.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, testCase.expectedResult, resp)
			}
		})
	}
}
