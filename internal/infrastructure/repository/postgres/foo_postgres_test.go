package postgres

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFooPostgres_FindAll(t *testing.T) {
	testCases := []struct {
		name          string
		pagination    handler.PaginationInput
		expectedCount int
		expectedError error
	}{
		{
			name:          "Success Case - Multiple Foos",
			pagination:    handler.PaginationInput{Offset: 0, Limit: 20},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:          "Success Case - Empty Foos",
			pagination:    handler.PaginationInput{Offset: 0, Limit: 0},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:          "Success Case - With Offset",
			pagination:    handler.PaginationInput{Offset: 1, Limit: 20},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name:          "Success Case - With Limit",
			pagination:    handler.PaginationInput{Offset: 0, Limit: 2},
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

func TestFooPostgres_FindByID(t *testing.T) {
	testCases := []struct {
		name           string
		id             int
		expectedResult model.Foo
		expectedError  error
	}{
		{
			name: "Success Case",
			id:   1,
			expectedResult: model.Foo{
				Id:     1,
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedError: nil,
		},
		{
			name:          "Fail Case - Not exist",
			id:            -1,
			expectedError: fmt.Errorf("foo with id 'id: -1' not found"),
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

			result, err := repo.FindByID(context.Background(), testCase.id)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, *result)
			}
		})
	}
}

func TestFooPostgres_Create(t *testing.T) {
	testCases := []struct {
		name          string
		foo           handler.FooCreateInput
		expectedError error
	}{
		{
			name: "Success Case",
			foo: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
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

			err := repo.Create(context.Background(), testCase.foo)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFooPostgres_Update(t *testing.T) {
	testCases := []struct {
		name          string
		foo           handler.FooUpdateInput
		expectedError error
	}{
		{
			name: "Success Case",
			foo: handler.FooUpdateInput{
				Id:     1,
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedError: nil,
		},
		{
			name: "Fail Case - Not exist",
			foo: handler.FooUpdateInput{
				Id:     -1,
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedError: fmt.Errorf("no row affected"),
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

			err := repo.Update(context.Background(), testCase.foo)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFooPostgres_DeleteByID(t *testing.T) {
	testCases := []struct {
		name          string
		id            int
		expectedError error
	}{
		{
			name:          "Success Case",
			id:            1,
			expectedError: nil,
		},
		{
			name:          "Fail Case - Not exist",
			id:            -1,
			expectedError: fmt.Errorf("no row affected"),
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

			err := repo.DeleteByID(context.Background(), testCase.id)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
