package redis

import (
	"astigo/internal/domain/model"
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationFooRedis_GetByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		id            uuid.UUID
		expectedError error
		expectedData  *model.Foo
	}{
		{
			name:          "Success Case",
			id:            uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedError: nil,
			expectedData: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
				Value:  1,
				Weight: 1.0,
			},
		},
		{
			name:          "Success Case - Not exist",
			id:            uuid.MustParse("40400000-0000-0000-0000-000000000000"),
			expectedError: nil,
			expectedData:  nil,
		},
	}

	ctx := context.Background()
	container, err := CreateRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	redis, err := NewRedis(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := NewFooRedis(redis)

			result, err := cache.GetByID(context.Background(), testCase.id)

			opts := []cmp.Option{
				cmpopts.IgnoreFields(model.Foo{}, "CreatedAt", "UpdatedAt"),
			}

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if testCase.expectedData == nil {
					assert.Nil(t, result)
				} else {
					assert.True(t, cmp.Equal(testCase.expectedData, result, opts...), cmp.Diff(testCase.expectedData, result, opts...))
				}
			}
		})
	}
}

func TestIntegrationFooRedis_Set(t *testing.T) {
	t.Parallel()
	now := time.Now()
	createdAt := now.Add(-1 * time.Hour)
	testCases := []struct {
		name          string
		foo           *model.Foo
		expectedError error
	}{
		{
			name: "Success Case - Create",
			foo: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000003"),
				Label:     "foo_created",
				Secret:    "secret_created",
				Value:     10,
				Weight:    1.5,
				CreatedAt: createdAt,
				UpdatedAt: nil,
			},
			expectedError: nil,
		},
		{
			name: "Success Case - Update",
			foo: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000003"),
				Label:     "foo_updated",
				Secret:    "secret_updated",
				Value:     20,
				Weight:    2.5,
				CreatedAt: createdAt,
				UpdatedAt: &now,
			},
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreateRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	redis, err := NewRedis(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := NewFooRedis(redis)

			err := cache.Set(ctx, testCase.foo, 0)

			opts := []cmp.Option{
				cmpopts.IgnoreFields(model.Foo{}, "CreatedAt", "UpdatedAt"),
			}

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				result, err := cache.GetByID(ctx, testCase.foo.Id)
				assert.NoError(t, err)
				assert.True(t, cmp.Equal(testCase.foo, result, opts...), cmp.Diff(testCase.foo, result, opts...))
			}
		})
	}
}

func TestIntegrationFooRedis_DeleteByID(t *testing.T) {
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
	}
	ctx := context.Background()
	container, err := CreateRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	redis, err := NewRedis(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := NewFooRedis(redis)

			err := cache.DeleteByID(ctx, testCase.id)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				result, err := cache.GetByID(ctx, testCase.id)
				assert.NoError(t, err)
				assert.Nil(t, result)
			}
		})
	}
}
