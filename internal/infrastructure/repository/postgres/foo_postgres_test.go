package postgres

import (
	"astigo/pkg/dto"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFooPostgres_FindAll(t *testing.T) {
	testCases := []struct {
		name          string
		pagination    dto.PaginationRequestDto
		expectedCount int
		expectedError error
	}{
		{
			name:          "Success Case - Multiple Foos",
			pagination:    dto.PaginationRequestDto{Offset: 0, Limit: 20},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:          "Success Case - Empty Foos",
			pagination:    dto.PaginationRequestDto{Offset: 0, Limit: 0},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:          "Success Case - With Offset",
			pagination:    dto.PaginationRequestDto{Offset: 1, Limit: 20},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name:          "Success Case - With Limit",
			pagination:    dto.PaginationRequestDto{Offset: 0, Limit: 2},
			expectedCount: 2,
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreatePostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	pg, err := NewPostgres(container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := NewFooPostgres(pg)

			result, err := repo.FindAll(context.Background(), testCase.pagination)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, testCase.expectedCount)
			}
		})
	}
}
