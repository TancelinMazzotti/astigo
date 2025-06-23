package redis

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationFooRedis_GetByID(t *testing.T) {
	testCases := []struct {
		name           string
		id             uuid.UUID
		expectedResult *model.Foo
		expectedError  error
	}{
		{
			name: "Success Case",
			id:   uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedResult: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedError: nil,
		},
		{
			name:           "Success Case - Not exist",
			id:             uuid.MustParse("40400000-0000-0000-0000-000000000000"),
			expectedResult: nil,
			expectedError:  nil,
		},
	}

	ctx := context.Background()
	container, err := CreateRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	redis, err := NewRedis(container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := NewFooRedis(redis)

			result, err := cache.GetByID(context.Background(), testCase.id)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if testCase.expectedResult == nil {
					assert.Nil(t, result)
					return
				} else {
					assert.Equal(t, *testCase.expectedResult, *result)
				}
			}
		})
	}
}

func TestIntegrationFooRedis_Set(t *testing.T) {
	testCases := []struct {
		name          string
		foo           model.Foo
		expectedError error
	}{
		{
			name: "Success Case - Create",
			foo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000003"),
				Label:  "foo_created",
				Secret: "secret_created",
			},
			expectedError: nil,
		},
		{
			name: "Success Case - Update",
			foo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000003"),
				Label:  "foo_updated",
				Secret: "secret_updated",
			},
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreateRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	redis, err := NewRedis(container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := NewFooRedis(redis)
			err := cache.Set(ctx, testCase.foo, 0)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				result, err := cache.GetByID(ctx, testCase.foo.Id)
				assert.NoError(t, err)
				assert.Equal(t, testCase.foo, *result)
			}
		})
	}
}

func TestIntegrationFooRedis_DeleteByID(t *testing.T) {
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

	redis, err := NewRedis(container.Config)
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
