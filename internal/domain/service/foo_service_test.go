package service

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func TestFooService_GetAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		input         handler.FooReadListInput
		expectedCount int
		expectedError error

		setupMockRepository func(*repository.MockFooRepository)
		setupMockCache      func(*cache.MockFooCache)
		setupMockMessaging  func(*messaging.MockFooMessaging)
	}{
		{
			name: "Success Case - Multiple Foos",
			input: handler.FooReadListInput{
				Offset: 0,
				Limit:  10,
			},
			expectedCount: 3,
			expectedError: nil,
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("FindAll", mock.Anything, handler.FooReadListInput{
					Offset: 0,
					Limit:  10,
				}).Return([]model.Foo{
					{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "Foo1"},
					{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "Foo2"},
					{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "Foo3"},
				}, nil)
			},
			setupMockCache:     func(mockCache *cache.MockFooCache) {},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
		{
			name: "Success Case - Empty Foos",
			input: handler.FooReadListInput{
				Offset: 0,
				Limit:  10,
			},
			expectedCount: 0,
			expectedError: nil,
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("FindAll", mock.Anything, handler.FooReadListInput{
					Offset: 0,
					Limit:  10,
				}).Return([]model.Foo{}, nil)
			},
			setupMockCache:     func(mockCache *cache.MockFooCache) {},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
		{
			name: "Failure Case - Repository Error",
			input: handler.FooReadListInput{
				Offset: 0,
				Limit:  10,
			},
			expectedCount: 0,
			expectedError: errors.New("fail to find all foo: repository error"),
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("FindAll", mock.Anything, handler.FooReadListInput{
					Offset: 0,
					Limit:  10,
				}).Return(([]model.Foo)(nil), errors.New("repository error"))
			},
			setupMockCache:     func(mockCache *cache.MockFooCache) {},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(repository.MockFooRepository)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)
			service := NewFooService(zap.NewNop(), mockRepo, mockCache, mockMessaging)

			testCase.setupMockRepository(mockRepo)
			testCase.setupMockCache(mockCache)
			testCase.setupMockMessaging(mockMessaging)

			result, err := service.GetAll(context.Background(), testCase.input)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, testCase.expectedCount)
			}
		})
	}
}

func TestFooService_GetByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		id             uuid.UUID
		expectedResult *model.Foo
		expectedError  error

		setupMockCache      func(*cache.MockFooCache)
		setupMockRepository func(*repository.MockFooRepository)
		setupMockMessaging  func(*messaging.MockFooMessaging)
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

			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return((*model.Foo)(nil), nil)

				mockCache.On("Set",
					mock.Anything,
					model.Foo{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "foo1",
						Secret: "secret1",
					}, FooCacheExpiration).Return(nil)
			},
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"FindByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(&model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo1",
					Secret: "secret1",
				}, nil)
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
		{
			name: "Success Case - Cache Hit",
			id:   uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedResult: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedError: nil,

			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(&model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo1",
					Secret: "secret1",
				}, nil)
			},
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {},
			setupMockMessaging:  func(mockMess *messaging.MockFooMessaging) {},
		},
		{
			name: "Success Case - Cache Error",
			id:   uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedResult: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo1",
				Secret: "secret1",
			},
			expectedError: nil,

			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return((*model.Foo)(nil), fmt.Errorf("cache error"))

				mockCache.On("Set",
					mock.Anything,
					model.Foo{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "foo1",
						Secret: "secret1",
					}, FooCacheExpiration).Return(fmt.Errorf("cache error"))
			},
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"FindByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(&model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo1",
					Secret: "secret1",
				}, nil)
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
		{
			name:           "Failure Case - Repository Error",
			id:             uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedResult: nil,
			expectedError:  fmt.Errorf("fail to find foo by id: repository error"),

			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return((*model.Foo)(nil), nil)

			},
			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"FindByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return((*model.Foo)(nil), fmt.Errorf("repository error"))
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(repository.MockFooRepository)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)
			service := NewFooService(zap.NewNop(), mockRepo, mockCache, mockMessaging)

			testCase.setupMockRepository(mockRepo)
			testCase.setupMockCache(mockCache)
			testCase.setupMockMessaging(mockMessaging)

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
	t.Parallel()
	testCases := []struct {
		name           string
		input          handler.FooCreateInput
		expectedResult *model.Foo
		expectedError  error

		setupMockRepository func(*repository.MockFooRepository)
		setupMockCache      func(*cache.MockFooCache)
		setupMockMessaging  func(*messaging.MockFooMessaging)
	}{
		{
			name: "Success Case",
			input: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedResult: &model.Foo{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedError: nil,

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On("Set", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				}), FooCacheExpiration).Return(nil)
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On("PublishFooCreated", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(nil)
			},
		},
		{
			name: "Success Case - Cache Error",
			input: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedResult: &model.Foo{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedError: nil,

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On("Set", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				}), FooCacheExpiration).Return(fmt.Errorf("cache error"))
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On("PublishFooCreated", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(nil)
			},
		},
		{
			name: "Failure Case - Repository Error",
			input: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedError: errors.New("fail to create foo: repository error"),

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(fmt.Errorf("repository error"))
			},
			setupMockCache:     func(mockCache *cache.MockFooCache) {},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
		{
			name: "Failure Case - Messaging Error",
			input: handler.FooCreateInput{
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedError: errors.New("fail to publish foo created: messaging error"),

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On("Set", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				}), FooCacheExpiration).Return(nil)
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On("PublishFooCreated", mock.Anything, mock.MatchedBy(func(foo model.Foo) bool {
					return foo.Label == "foo_create" &&
						foo.Secret == "secret_create"
				})).Return(fmt.Errorf("messaging error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(repository.MockFooRepository)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)
			service := NewFooService(zap.NewNop(), mockRepo, mockCache, mockMessaging)

			testCase.setupMockRepository(mockRepo)
			testCase.setupMockCache(mockCache)
			testCase.setupMockMessaging(mockMessaging)

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
	t.Parallel()
	testCases := []struct {
		name          string
		input         handler.FooUpdateInput
		expectedError error

		setupMockRepository func(*repository.MockFooRepository)
		setupMockCache      func(*cache.MockFooCache)
		setupMockMessaging  func(*messaging.MockFooMessaging)
	}{
		{
			name: "Success Case",
			input: handler.FooUpdateInput{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedError: nil,

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"FindByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(&model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo1",
					Secret: "secret1",
				}, nil)

				mockRepo.On("Update", mock.Anything, model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On("Set", mock.Anything, model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}, FooCacheExpiration).Return(nil)
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On("PublishFooUpdated", mock.Anything, model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}).Return(nil)
			},
		},
		{
			name: "Success Case - Cache Error",
			input: handler.FooUpdateInput{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedError: nil,

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"FindByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(&model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo1",
					Secret: "secret1",
				}, nil)

				mockRepo.On("Update", mock.Anything, model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On("Set", mock.Anything, model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}, FooCacheExpiration).Return(fmt.Errorf("cache error"))
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On("PublishFooUpdated", mock.Anything, model.Foo{
					Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:  "foo_update",
					Secret: "secret_update",
				}).Return(nil)
			},
		},
		{
			name: "Failure Case - Repository Get Error",
			input: handler.FooUpdateInput{
				Id:     uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				Label:  "foo_update",
				Secret: "secret_update",
			},
			expectedError: errors.New("fail to get foo by id: repository error"),

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"FindByID",
					mock.Anything,
					uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				).Return((*model.Foo)(nil), errors.New("repository error"))
			},
			setupMockCache:     func(mockCache *cache.MockFooCache) {},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(repository.MockFooRepository)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)
			service := NewFooService(zap.NewNop(), mockRepo, mockCache, mockMessaging)

			testCase.setupMockRepository(mockRepo)
			testCase.setupMockCache(mockCache)
			testCase.setupMockMessaging(mockMessaging)

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
	t.Parallel()
	testCases := []struct {
		name          string
		id            uuid.UUID
		expectedError error

		setupMockRepository func(*repository.MockFooRepository)
		setupMockCache      func(*cache.MockFooCache)
		setupMockMessaging  func(*messaging.MockFooMessaging)
	}{
		{
			name:          "Success Case",
			id:            uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedError: nil,

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On(
					"PublishFooDeleted",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
		},
		{
			name:          "Success Case - Cache Error",
			id:            uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedError: nil,

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
			setupMockCache: func(mockCache *cache.MockFooCache) {
				mockCache.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(fmt.Errorf("cache error"))
			},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {
				mockMess.On(
					"PublishFooDeleted",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
		},
		{
			name:          "Failure Case - Repository Error",
			id:            uuid.MustParse("40000000-0000-0000-0000-000000000000"),
			expectedError: errors.New("fail to delete foo by id: repository error"),

			setupMockRepository: func(mockRepo *repository.MockFooRepository) {
				mockRepo.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				).Return(errors.New("repository error"))
			},
			setupMockCache:     func(mockCache *cache.MockFooCache) {},
			setupMockMessaging: func(mockMess *messaging.MockFooMessaging) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(repository.MockFooRepository)
			mockCache := new(cache.MockFooCache)
			mockMessaging := new(messaging.MockFooMessaging)
			service := NewFooService(zap.NewNop(), mockRepo, mockCache, mockMessaging)

			testCase.setupMockRepository(mockRepo)
			testCase.setupMockCache(mockCache)
			testCase.setupMockMessaging(mockMessaging)

			err := service.DeleteByID(context.Background(), testCase.id)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
