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

			service := NewService(mockRepo, mockCache, mockMessaging)

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
