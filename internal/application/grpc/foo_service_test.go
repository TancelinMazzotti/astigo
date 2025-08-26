package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/contract/data"
	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/mocks/domain/contract/service"
	"github.com/TancelinMazzotti/astigo/pkg/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFooService_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		request       *proto.ListFoosRequest
		expectedError error
		expectedCount int

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name: "Success Case",
			request: &proto.ListFoosRequest{
				Offset: 0,
				Limit:  10,
			},
			expectedCount: 3,
			expectedError: nil,

			setupMockHandler: func(mockRepo *service.MockFooService) {
				mockRepo.On("GetAll",
					mock.Anything,
					data.FooReadListInput{Offset: 0, Limit: 10},
				).Return(
					[]*model.Foo{
						{
							Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
							Label:     "foo1",
							Secret:    "secret1",
							Value:     1,
							Weight:    1.5,
							CreatedAt: time.Now(),
							UpdatedAt: nil,
						},
						{
							Id:        uuid.MustParse("20000000-0000-0000-0000-000000000002"),
							Label:     "foo2",
							Secret:    "secret2",
							Value:     2,
							Weight:    2.5,
							CreatedAt: time.Now(),
							UpdatedAt: nil,
						},
						{
							Id:        uuid.MustParse("20000000-0000-0000-0000-000000000003"),
							Label:     "foo3",
							Secret:    "secret3",
							Value:     3,
							Weight:    3.5,
							CreatedAt: time.Now(),
							UpdatedAt: nil,
						}},
					nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			svc := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

			resp, err := svc.List(context.Background(), testCase.request)

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

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name: "Success Case",
			request: &proto.GetFooRequest{
				Id: "20000000-0000-0000-0000-000000000001",
			},
			expectedError: nil,
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:     "20000000-0000-0000-0000-000000000001",
					Label:  "Foo1",
					Value:  1,
					Weight: 1.5,
				},
			},

			setupMockHandler: func(mockRepo *service.MockFooService) {
				mockRepo.On("GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(
					&model.Foo{
						Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:     "Foo1",
						Value:     1,
						Weight:    1.5,
						CreatedAt: time.Now(),
						UpdatedAt: nil,
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
			setupMockHandler: func(mockRepo *service.MockFooService) {},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			svc := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

			resp, err := svc.Get(context.Background(), testCase.request)

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

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name: "Success Case",
			request: &proto.CreateFooRequest{
				Label:  "foo_create",
				Secret: "secret_create",
				Value:  1,
				Weight: 1.5,
			},
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:     "20000000-0000-0000-0000-000000000001",
					Label:  "foo_create",
					Value:  1,
					Weight: 1.5,
				},
			},
			expectedError: nil,

			setupMockHandler: func(mockRepo *service.MockFooService) {
				mockRepo.On("Create",
					mock.Anything,
					data.FooCreateInput{
						Label:  "foo_create",
						Secret: "secret_create",
						Value:  1,
						Weight: 1.5,
					},
				).Return(&model.Foo{
					Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:     "foo_create",
					Secret:    "secret_create",
					Value:     1,
					Weight:    1.5,
					CreatedAt: time.Now(),
					UpdatedAt: nil,
				}, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			svc := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

			resp, err := svc.Create(context.Background(), testCase.request)

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

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name: "Success Case",
			request: &proto.UpdateFooRequest{
				Id:     "20000000-0000-0000-0000-000000000001",
				Label:  "foo_update",
				Secret: "secret_update",
				Value:  1,
				Weight: 1.5,
			},
			expectedResult: &proto.FooResponse{
				Foo: &proto.Foo{
					Id:     "20000000-0000-0000-0000-000000000001",
					Label:  "foo_update",
					Value:  1,
					Weight: 1.5,
				},
			},
			expectedError: nil,

			setupMockHandler: func(mockRepo *service.MockFooService) {
				mockRepo.On("Update", mock.Anything, &data.FooUpdateInput{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
					Value:  1,
					Weight: 1.5,
				}).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			svc := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

			resp, err := svc.Update(context.Background(), testCase.request)

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

		setupMockHandler func(*service.MockFooService)
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

			setupMockHandler: func(mockRepo *service.MockFooService) {
				mockRepo.On("DeleteByID", mock.Anything, uuid.MustParse("20000000-0000-0000-0000-000000000001")).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			svc := NewFooService(mockHandler)

			testCase.setupMockHandler(mockHandler)

			resp, err := svc.Delete(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
				assert.Nil(t, resp)
			}
		})
	}
}
