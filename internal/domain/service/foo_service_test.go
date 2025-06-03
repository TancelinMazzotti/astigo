package service

import (
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
		name          string
		mockResponse  []dto.FooResponseReadDto
		mockError     error
		expectedCount int
		expectedError error
	}{
		{
			name: "Success Case - Multiple Foos",
			mockResponse: []dto.FooResponseReadDto{
				{Id: 1, Label: "Foo1"},
				{Id: 2, Label: "Foo2"},
			},
			mockError:     nil,
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name:          "Failure Case - Repository Error",
			mockResponse:  nil,
			mockError:     errors.New("repository error"),
			expectedCount: 0,
			expectedError: errors.New("repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(testCase.mockResponse, testCase.mockError)

			service := FooService{repo: mockRepo}

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
