package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"errors"
	"github.com/google/uuid"
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
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "Foo1"},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "Foo2"},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "Foo3"},
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
		name              string
		id                uuid.UUID
		mockCacheResponse *model.Foo
		mockCacheError    error
		isCached          bool
		mockRepoResponse  *model.Foo
		mockRepoError     error
		expectedResult    *model.Foo
		expectedError     error
	}{
		{
			name:              "Success Case",
			id:                uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockCacheResponse: nil,
			mockCacheError:    nil,
			isCached:          false,
			mockRepoResponse: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedResult: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			mockRepoError: nil,
		},
		{
			name: "Success Case - Cache Hit",
			id:   uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockCacheResponse: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			mockCacheError: nil,
			isCached:       true,
			expectedResult: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			mockRepoError: nil,
		},
		{
			name:              "Failure Case - Cache Error",
			id:                uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockCacheResponse: nil,
			mockCacheError:    errors.New("cache error"),
			expectedError:     errors.New("fail to find foo by id from cache: cache error"),
		},
		{
			name:              "Failure Case - Repository Error",
			id:                uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockCacheResponse: nil,
			mockCacheError:    nil,
			isCached:          false,
			mockRepoResponse:  nil,
			mockRepoError:     errors.New("repository error"),
			expectedError:     errors.New("fail to find foo by id: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockCache := new(cache.MockFooCache)
			mockCache.On("GetByID", mock.Anything, testCase.id).Return(testCase.mockCacheResponse, testCase.mockCacheError)

			mockRepo := new(repository.MockFooRepository)
			if !testCase.isCached {
				mockRepo.On("FindByID", mock.Anything, testCase.id).Return(testCase.mockRepoResponse, testCase.mockRepoError)
				if (testCase.mockRepoResponse != nil) && (testCase.mockRepoError == nil) {
					mockCache.On("Set", mock.Anything, *testCase.mockRepoResponse, FooCacheExpiration).Return(nil)
				}
			}
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
		name            string
		input           handler.FooCreateInput
		mockRepoRequest model.Foo
		mockRepoError   error
		expectedResult  *model.Foo
		expectedError   error
	}{
		{
			name: "Success Case",
			input: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockRepoRequest: model.Foo{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockRepoError: nil,
			expectedResult: &model.Foo{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedError: nil,
		},
		{
			name: "Failure Case - Repository Error",
			input: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			mockRepoRequest: model.Foo{
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
			mockRepo.On("Create", mock.Anything, mock.Anything).Return(testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			foo, err := service.Create(context.Background(), testCase.input)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, foo.Id)
				assert.Equal(t, testCase.expectedResult.Label, foo.Label)
				assert.Equal(t, testCase.expectedResult.Secret, foo.Secret)
			}
		})
	}
}

func TestFooService_Update(t *testing.T) {
	testCases := []struct {
		name  string
		input handler.FooUpdateInput

		mockRepoGetResponse   *model.Foo
		mockRepoGetError      error
		mockRepoUpdateRequest model.Foo
		mockRepoUpdateError   error

		mockCacheRequest model.Foo
		mockCacheError   error
		expectedError    error
	}{
		{
			name: "Success Case",
			input: handler.FooUpdateInput{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},

			mockRepoGetResponse: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},

			mockRepoUpdateRequest: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},

			mockCacheRequest: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
		},
		{
			name: "Failure Case - Repository Error",
			input: handler.FooUpdateInput{
				Id:     uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				Label:  "foo_update",
				Secret: "secret_update",
			},

			mockRepoGetResponse: &model.Foo{
				Id:     uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				Label:  "foo1",
				Secret: "secret1",
			},

			mockRepoUpdateRequest: model.Foo{
				Id:     uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
			mockRepoUpdateError: errors.New("repository error"),
			expectedError:       errors.New("fail to update foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("FindByID", mock.Anything, testCase.input.Id).Return(testCase.mockRepoGetResponse, testCase.mockRepoGetError)
			mockRepo.On("Update", mock.Anything, testCase.mockRepoUpdateRequest).Return(testCase.mockRepoUpdateError)

			mockCache := new(cache.MockFooCache)
			mockCache.On("Set", mock.Anything, testCase.mockCacheRequest, FooCacheExpiration).Return(testCase.mockCacheError)
			mockMessaging := new(messaging.MockFooMessaging)

			service := NewFooService(mockRepo, mockCache, mockMessaging)

			err := service.Update(context.Background(), testCase.input)

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
		name           string
		id             uuid.UUID
		mockRepoError  error
		mockCacheError error
		isCached       bool
		expectedError  error
	}{
		{
			name:           "Success Case",
			id:             uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockRepoError:  nil,
			mockCacheError: nil,
			isCached:       true,
			expectedError:  nil,
		},
		{
			name:           "Failure Case - Repository Error",
			id:             uuid.MustParse("40000000-0000-0000-0000-000000000000"),
			mockRepoError:  errors.New("repository error"),
			mockCacheError: nil,
			isCached:       false,
			expectedError:  errors.New("fail to delete foo by id: repository error"),
		},
		{
			name:           "Failure Case - Cache Error",
			id:             uuid.MustParse("40000000-0000-0000-0000-000000000000"),
			mockRepoError:  nil,
			mockCacheError: errors.New("cache error"),
			isCached:       true,
			expectedError:  errors.New("fail to delete foo by id from cache: cache error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := new(repository.MockFooRepository)
			mockRepo.On("DeleteByID", mock.Anything, testCase.id).Return(testCase.mockRepoError)
			mockCache := new(cache.MockFooCache)
			if testCase.isCached {
				mockCache.On("DeleteByID", mock.Anything, testCase.id).Return(testCase.mockCacheError)
			}
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
