package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/repository"
	"astigo/pkg/dto"
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
		mockRepoResponse []dto.FooResponseReadDto
		mockRepoError    error
		expectedCount    int
		expectedError    error
	}{
		{
			name: "Success Case - Multiple Foos",
			mockRepoResponse: []dto.FooResponseReadDto{
				{Id: 1, Label: "Foo1", Bars: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
				{Id: 2, Label: "Foo2", Bars: []int{}},
				{Id: 3, Label: "Foo3", Bars: []int{}},
			},
			mockRepoError: nil,
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:             "Success Case - Empty Foos",
			mockRepoResponse: []dto.FooResponseReadDto{},
			mockRepoError:    nil,
			expectedCount:    0,
			expectedError:    nil,
		},
		{
			name:             "Failure Case - Repository Error",
			mockRepoResponse: nil,
			mockRepoError:    errors.New("repository error"),
			expectedCount:    0,
			expectedError:    errors.New("fail to find all foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(testCase.mockRepoResponse, testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewService(mockRepo, mockCache, mockMessaging)

			result, err := service.GetAll(context.Background(), dto.PaginationRequestDto{})

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
