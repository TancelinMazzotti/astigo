package postgres

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestIntegrationFooPostgres_FindAll verifies the behavior of the FindAll repository method with Postgres integration.
func TestIntegrationFooPostgres_FindAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		input         handler.FooReadListInput
		expectedCount int
		expectedError error
		expectedData  []*model.Foo
	}{
		{
			name:          "Success Case - Multiple Foos",
			input:         handler.FooReadListInput{Offset: 0, Limit: 20},
			expectedCount: 3,
			expectedError: nil,
			expectedData: []*model.Foo{
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "foo1", Secret: "secret1", Value: 1, Weight: 1.0, UpdatedAt: nil},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "foo2", Secret: "secret2", Value: 2, Weight: 2.0, UpdatedAt: nil},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "foo3", Secret: "secret3", Value: 3, Weight: 3.0, UpdatedAt: nil},
			},
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
			expectedData: []*model.Foo{
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "foo2", Secret: "secret2", Value: 2, Weight: 2.0, UpdatedAt: nil},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "foo3", Secret: "secret3", Value: 3, Weight: 3.0, UpdatedAt: nil},
			},
		},
		{
			name:          "Success Case - With Limit",
			input:         handler.FooReadListInput{Offset: 0, Limit: 2},
			expectedCount: 2,
			expectedError: nil,
			expectedData: []*model.Foo{
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "foo1", Secret: "secret1", Value: 1, Weight: 1.0, UpdatedAt: nil},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "foo2", Secret: "secret2", Value: 2, Weight: 2.0, UpdatedAt: nil},
			},
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

			// Ignore the CreatedAt field in comparison since it's automatically set by the database at insertion time
			opts := []cmp.Option{
				cmpopts.IgnoreFields(model.Foo{}, "CreatedAt"),
			}

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, testCase.expectedCount)
				assert.True(t, cmp.Equal(testCase.expectedData, result, opts...), cmp.Diff(testCase.expectedData, result, opts...))

				for i := range result {
					assert.NotZero(t, result[i].CreatedAt)
				}

			}
		})
	}
}

// TestIntegrationFooPostgres_FindByID tests the integration of the FindByID method for the FooPostgres repository with a PostgreSQL database.
func TestIntegrationFooPostgres_FindByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		id            uuid.UUID
		expectedData  *model.Foo
		expectedError error
	}{
		{
			name: "Success Case",
			id:   uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedData: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "foo1",
				Secret:    "secret1",
				Value:     1,
				Weight:    1.0,
				UpdatedAt: nil,
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

			// Ignore the CreatedAt field in comparison since it's automatically set by the database at insertion time
			opts := []cmp.Option{
				cmpopts.IgnoreFields(model.Foo{}, "CreatedAt"),
			}

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.True(t, cmp.Equal(testCase.expectedData, result, opts...), cmp.Diff(testCase.expectedData, result, opts...))
				assert.NotZero(t, result.CreatedAt)
			}
		})
	}
}

// TestIntegrationFooPostgres_Create validates the creation of a Foo record in the Postgres repository under integration tests.
func TestIntegrationFooPostgres_Create(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		foo           *model.Foo
		expectedError error
	}{
		{
			name: "Success Case",
			foo: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-100000000000"),
				Label:  "foo_create",
				Secret: "secret_create",
				Value:  50,
				Weight: 1.5,
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

// TestIntegrationFooPostgres_Update tests the update functionality of the FooPostgres repository with integration against Postgres.
func TestIntegrationFooPostgres_Update(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		foo           *model.Foo
		expectedError error
	}{
		{
			name: "Success Case",
			foo: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
				Value:  50,
				Weight: 1.5,
			},
			expectedError: nil,
		},
		{
			name: "Fail Case - Not exist",
			foo: &model.Foo{
				Id:     uuid.MustParse("40400000-0000-0000-0000-000000000000"),
				Label:  "foo_update",
				Secret: "secret_update",
				Value:  50,
				Weight: 1.5,
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
				assert.NotNil(t, testCase.foo.UpdatedAt)
				assert.NotZero(t, testCase.foo.UpdatedAt)
			}
		})
	}
}

// TestIntegrationFooPostgres_DeleteByID tests the DeleteByID function of FooPostgres in an integration Postgres database scenario.
// It validates both successful deletion and failure when attempting to delete a non-existent record.
func TestIntegrationFooPostgres_DeleteByID(t *testing.T) {
	t.Parallel()
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
