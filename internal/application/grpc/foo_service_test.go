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
	t.Parallel()
	testCases := []struct {
		name          string
		request       *proto.ListFoosRequest
		expectedError error
		expectedCount int

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name: "Success Case",
			request: &proto.ListFoosRequest{
				Offset: 0,
				Limit:  10,
			},
			expectedCount: 3,
			expectedError: nil,

			setupMockHandler: func(mockRepo *handler.MockFooHandler) {
				mockRepo.On("GetAll",
					mock.Anything,
					handler.FooReadListInput{Offset: 0, Limit: 10},
				).Return(
					[]model.Foo{
						{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "Foo1"},
						{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "Foo2"},
						{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "Foo3"}},
					nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
			service := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

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
		})
	}
}

func TestFooService_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		request        *proto.GetFooRequest
		expectedError  error
		expectedResult *proto.FooResponse

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name: "Success Case",
			request: &proto.GetFooRequest{
				Id: "20000000-0000-0000-0000-000000000001",
			},
			expectedError: nil,
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:    "20000000-0000-0000-0000-000000000001",
					Label: "Foo1",
				},
			},

			setupMockHandler: func(mockRepo *handler.MockFooHandler) {
				mockRepo.On("GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(
					&model.Foo{
						Id:    uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label: "Foo1",
					}, nil)
			},
		},
		{
			name: "Failed Case - Not UUID",
			request: &proto.GetFooRequest{
				Id: "not uuid",
			},
			expectedError:    fmt.Errorf("fail to parse id"),
			expectedResult:   nil,
			setupMockHandler: func(mockRepo *handler.MockFooHandler) {},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
			service := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

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
	t.Parallel()
	testCases := []struct {
		name           string
		request        *proto.CreateFooRequest
		expectedResult *proto.FooResponse
		expectedError  error

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name: "Success Case",
			request: &proto.CreateFooRequest{
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

			setupMockHandler: func(mockRepo *handler.MockFooHandler) {
				mockRepo.On("Create",
					mock.Anything,
					handler.FooCreateInput{
						Label:  "foo_create",
						Secret: "secret_create",
					},
				).Return(&model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_create",
					Secret: "secret_create",
				}, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
			service := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

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
	t.Parallel()
	testCases := []struct {
		name           string
		request        *proto.UpdateFooRequest
		expectedResult *proto.FooResponse
		expectedError  error

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name: "Success Case",
			request: &proto.UpdateFooRequest{
				Id:     "20000000-0000-0000-0000-000000000001",
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:    "20000000-0000-0000-0000-000000000001",
					Label: "foo_update",
				},
			},
			expectedError: nil,

			setupMockHandler: func(mockRepo *handler.MockFooHandler) {
				mockRepo.On("Update", mock.Anything, handler.FooUpdateInput{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
			service := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

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
	t.Parallel()
	testCases := []struct {
		name           string
		request        *proto.DeleteFooRequest
		expectedResult *proto.DeleteFooResponse
		expectedError  error

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name: "Success Case",
			request: &proto.DeleteFooRequest{
				Id: "20000000-0000-0000-0000-000000000001",
			},
			expectedResult: &proto.DeleteFooResponse{
				Success: true,
			},
			expectedError: nil,

			setupMockHandler: func(mockRepo *handler.MockFooHandler) {
				mockRepo.On("DeleteByID", mock.Anything, uuid.MustParse("20000000-0000-0000-0000-000000000001")).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
			service := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

			resp, err := service.Delete(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			}
		})
	}
}
