package grpc

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"astigo/pkg/proto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFooService_List(t *testing.T) {
	testCases := []struct {
		name    string
		request *proto.ListFoosRequest

		mockRequest  handler.FooReadListInput
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

			mockRequest: handler.FooReadListInput{Offset: 0, Limit: 10},
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
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Foos, testCase.expectedCount)
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
		{
			name: "Failed Case - Not UUID",
			request: &proto.GetFooRequest{
				Id: "not uuid",
			},
			expectedError: fmt.Errorf("fail to parse id"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("GetByID", mock.Anything, testCase.mockRequest).Return(testCase.mockResponse, testCase.mockError)

			service := NewFooService(mockHandler)

			resp, err := service.Get(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, testCase.expectedResult, resp)
			}
		})
	}
}

func TestFooService_Create(t *testing.T) {
	testCases := []struct {
		name           string
		request        *proto.CreateFooRequest
		mockRequest    handler.FooCreateInput
		mockResponse   *model.Foo
		mockError      error
		expectedResult *proto.FooResponse
		expectedError  error
	}{
		{
			name: "Success Case",
			request: &proto.CreateFooRequest{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockRequest: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockResponse: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_create",
				Secret: "secret_create",
			},

			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:    "20000000-0000-0000-0000-000000000001",
					Label: "foo_create",
				},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("Create", mock.Anything, testCase.mockRequest).Return(testCase.mockResponse, testCase.mockError)

			service := NewFooService(mockHandler)

			resp, err := service.Create(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, testCase.expectedResult, resp)
			}
		})
	}
}

func TestFooService_Update(t *testing.T) {
	testCases := []struct {
		name           string
		request        *proto.UpdateFooRequest
		mockRequest    handler.FooUpdateInput
		mockError      error
		expectedResult *proto.FooResponse
		expectedError  error
	}{
		{
			name: "Success Case",
			request: &proto.UpdateFooRequest{
				Id:     "20000000-0000-0000-0000-000000000001",
				Label:  "foo_update",
				Secret: "secret_update",
			},
			mockRequest: handler.FooUpdateInput{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
			mockError: nil,
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:    "20000000-0000-0000-0000-000000000001",
					Label: "foo_update",
				},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("Update", mock.Anything, testCase.mockRequest).Return(testCase.mockError)

			service := NewFooService(mockHandler)

			resp, err := service.Update(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			}
		})
	}
}

func TestFooService_Delete(t *testing.T) {
	testCases := []struct {
		name           string
		request        *proto.DeleteFooRequest
		mockRequest    uuid.UUID
		mockError      error
		expectedResult *proto.DeleteFooResponse
		expectedError  error
	}{
		{
			name: "Success Case",
			request: &proto.DeleteFooRequest{
				Id: "20000000-0000-0000-0000-000000000001",
			},
			mockRequest: uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockError:   nil,
			expectedResult: &proto.DeleteFooResponse{
				Success: true,
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("DeleteByID", mock.Anything, testCase.mockRequest).Return(testCase.mockError)

			service := NewFooService(mockHandler)

			resp, err := service.Delete(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			}
		})
	}
}
