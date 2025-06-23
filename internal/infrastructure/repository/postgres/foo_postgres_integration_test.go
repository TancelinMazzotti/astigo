package postgres

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationFooPostgres_FindAll(t *testing.T) {
	testCases := []struct {
		name          string
		input         handler.FooReadListInput
		expectedCount int
		expectedError error
	}{
		{
			name:          "Success Case - Multiple Foos",
			input:         handler.FooReadListInput{Offset: 0, Limit: 20},
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:          "Success Case - Empty Foos",
			input:         handler.FooReadListInput{Offset: 0, Limit: 0},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:          "Success Case - With Offset",
			input:         handler.FooReadListInput{Offset: 1, Limit: 20},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name:          "Success Case - With Limit",
			input:         handler.FooReadListInput{Offset: 0, Limit: 2},
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

			result, err := repo.FindAll(context.Background(), testCase.input)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, testCase.expectedCount)
			}
		})
	}
}

func TestIntegrationFooPostgres_FindByID(t *testing.T) {
	testCases := []struct {
		name           string
		id             uuid.UUID
		expectedResult model.Foo
		expectedError  error
	}{
		{
			name: "Success Case",
			id:   uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedResult: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedError: nil,
		},
		{
			name:          "Fail Case - Not exist",
			id:            uuid.MustParse("40400000-0000-0000-0000-000000000000"),
			expectedError: fmt.Errorf("foo with id 'id: 40400000-0000-0000-0000-000000000000' not found"),
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

func TestIntegrationFooPostgres_Create(t *testing.T) {
	testCases := []struct {
		name          string
		foo           model.Foo
		expectedError error
	}{
		{
			name: "Success Case",
			foo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-100000000000"),
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

func TestIntegrationFooPostgres_Update(t *testing.T) {
	testCases := []struct {
		name          string
		foo           model.Foo
		expectedError error
	}{
		{
			name: "Success Case",
			foo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedError: nil,
		},
		{
			name: "Fail Case - Not exist",
			foo: model.Foo{
				Id:     uuid.MustParse("40400000-0000-0000-0000-000000000000"),
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

func TestIntegrationFooPostgres_DeleteByID(t *testing.T) {
	testCases := []struct {
		name          string
		id            uuid.UUID
		expectedError error
	}{
		{
			name:          "Success Case",
			id:            uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedError: nil,
		},
		{
			name:          "Fail Case - Not exist",
			id:            uuid.MustParse("40400000-0000-0000-0000-000000000000"),
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
