package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFooService_GetAll(t *testing.T) {
	// Scénarios de test
	testCases := []struct {
		name             string
		mockRepoRequest  handler.PaginationInput
		mockRepoResponse []model.Foo
		mockRepoError    error
		expectedCount    int
		expectedError    error
	}{
		{
			name: "Success Case - Multiple Foos",
			mockRepoRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockRepoResponse: []model.Foo{
				{Id: 1, Label: "Foo1"},
				{Id: 2, Label: "Foo2"},
				{Id: 3, Label: "Foo3"},
			},
			mockRepoError: nil,
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name: "Success Case - Empty Foos",
			mockRepoRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockRepoResponse: []model.Foo{},
			mockRepoError:    nil,
			expectedCount:    0,
			expectedError:    nil,
		},
		{
			name: "Failure Case - Repository Error",
			mockRepoRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockRepoResponse: nil,
			mockRepoError:    errors.New("repository error"),
			expectedCount:    0,
			expectedError:    errors.New("fail to find all foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("FindAll", mock.Anything, testCase.mockRepoRequest).Return(testCase.mockRepoResponse, testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			result, err := service.GetAll(context.Background(), testCase.mockRepoRequest)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, testCase.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFooService_GetByID(t *testing.T) {
	testCases := []struct {
		name             string
		id               int
		mockRepoResponse *model.Foo
		mockRepoError    error
		expectedResult   *model.Foo
		expectedError    error
	}{
		{
			name: "Success Case",
			id:   1,
			mockRepoResponse: &model.Foo{
				Id:     1,
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedResult: &model.Foo{
				Id:     1,
				Label:  "foo1",
				Secret: "secret1",
			},
			mockRepoError: nil,
		},
		{
			name:             "Failure Case - Repository Error",
			id:               1,
			mockRepoResponse: nil,
			mockRepoError:    errors.New("repository error"),
			expectedError:    errors.New("fail to find foo by id: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("FindByID", mock.Anything, testCase.id).Return(testCase.mockRepoResponse, testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			result, err := service.GetByID(context.Background(), testCase.id)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, result)
			}
		})
	}
}

func TestFooService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		foo           handler.FooCreateInput
		mockRepoError error
		expectedError error
	}{
		{
			name: "Success Case",
			foo: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockRepoError: nil,
			expectedError: nil,
		},
		{
			name: "Failure Case - Repository Error",
			foo: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockRepoError: errors.New("repository error"),
			expectedError: errors.New("fail to create foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("Create", mock.Anything, testCase.foo).Return(testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			err := service.Create(context.Background(), testCase.foo)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFooService_Update(t *testing.T) {
	testCases := []struct {
		name          string
		foo           handler.FooUpdateInput
		mockRepoError error
		expectedError error
	}{
		{
			name: "Success Case",
			foo: handler.FooUpdateInput{
				Id:     1,
				Label:  "foo_update",
				Secret: "secret_update",
			},
			mockRepoError: nil,
			expectedError: nil,
		},
		{
			name: "Failure Case - Repository Error",
			foo: handler.FooUpdateInput{
				Id:     -1,
				Label:  "foo_update",
				Secret: "secret_update",
			},
			mockRepoError: errors.New("repository error"),
			expectedError: errors.New("fail to update foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("Update", mock.Anything, testCase.foo).Return(testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			err := service.Update(context.Background(), testCase.foo)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFooService_DeleteByID(t *testing.T) {
	testCases := []struct {
		name          string
		id            int
		mockRepoError error
		expectedError error
	}{
		{
			name:          "Success Case",
			id:            1,
			mockRepoError: nil,
			expectedError: nil,
		},
		{
			name:          "Failure Case - Repository Error",
			id:            -1,
			mockRepoError: errors.New("repository error"),
			expectedError: errors.New("fail to delete foo by id: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("DeleteByID", mock.Anything, testCase.id).Return(testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			err := service.DeleteByID(context.Background(), testCase.id)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
