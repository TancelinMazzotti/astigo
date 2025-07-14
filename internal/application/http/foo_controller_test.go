package http

import (
	"astigo/internal/domain/adapter/data"
	"astigo/internal/domain/adapter/repository"
	"astigo/internal/domain/model"
	"astigo/internal/domain/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFooController_GetAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		statusCode   int
		bodyResponse string

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name:       "Success Case - Multiple Foos",
			url:        "/foos?offset=0&limit=10",
			statusCode: http.StatusOK,
			bodyResponse: `[
				{"id":"20000000-0000-0000-0000-000000000001", "label":"foo1", "value":1, "weight":1.5},
				{"id":"20000000-0000-0000-0000-000000000002", "label":"foo2", "value":2, "weight":2.5},
				{"id":"20000000-0000-0000-0000-000000000003", "label":"foo3", "value":3, "weight":3.5}
			]`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"GetAll",
					mock.Anything,
					data.FooReadListInput{Offset: 0, Limit: 10},
				).Return([]*model.Foo{
					{
						Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:     "foo1",
						Secret:    "secret1",
						Value:     1,
						Weight:    1.5,
						CreatedAt: time.Now(),
					},
					{
						Id:        uuid.MustParse("20000000-0000-0000-0000-000000000002"),
						Label:     "foo2",
						Secret:    "secret2",
						Value:     2,
						Weight:    2.5,
						CreatedAt: time.Now(),
					},
					{
						Id:        uuid.MustParse("20000000-0000-0000-0000-000000000003"),
						Label:     "foo3",
						Secret:    "secret3",
						Value:     3,
						Weight:    3.5,
						CreatedAt: time.Now(),
					},
				}, nil)
			},
		},
		{
			name:         "Success Case - No Foos",
			url:          "/foos?offset=0&limit=10",
			statusCode:   http.StatusOK,
			bodyResponse: `[]`,
			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"GetAll",
					mock.Anything,
					data.FooReadListInput{Offset: 0, Limit: 10},
				).Return([]*model.Foo{}, nil)
			},
		},
		{
			name:             "Failure Case - Invalid type offset",
			url:              "/foos?offset=invalid&limit=10",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate query params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:             "Failure Case - Invalid value offset",
			url:              "/foos?offset=-1&limit=10",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate query params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:             "Failure Case - Invalid type limit",
			url:              "/foos?offset=0&limit=invalid",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate query params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:             "Failure Case - Invalid value limit",
			url:              "/foos?offset=0&limit=-1",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate query params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:             "Failure Case - Invalid exceeded limit",
			url:              "/foos?offset=0&limit=51",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate query params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:         "Failure Case - Repository Error",
			url:          "/foos?offset=0&limit=10",
			statusCode:   http.StatusInternalServerError,
			bodyResponse: `{"error":"failed to get all foos"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"GetAll",
					mock.Anything,
					data.FooReadListInput{Offset: 0, Limit: 10},
				).Return(([]*model.Foo)(nil), errors.New("repository error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			controller := NewFooController(mockHandler)

			testCase.setupMockHandler(mockHandler)

			req, err := http.NewRequest(http.MethodGet, testCase.url, nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.GET("/foos", controller.GetAll)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
			assert.JSONEq(t, testCase.bodyResponse, w.Body.String())
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestFooController_GetByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		statusCode   int
		bodyResponse string

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name:         "Success Case",
			url:          "/foos/20000000-0000-0000-0000-000000000001",
			statusCode:   http.StatusOK,
			bodyResponse: `{"id":"20000000-0000-0000-0000-000000000001", "label":"foo1", "value":1, "weight":1.5}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(
					&model.Foo{
						Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:     "foo1",
						Secret:    "secret1",
						Value:     1,
						Weight:    1.5,
						CreatedAt: time.Now(),
					}, nil)
			},
		},
		{
			name:             "Failure Case - Not UUID",
			url:              "/foos/not_uuid",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate path params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:         "Failure Case - Not Found",
			url:          "/foos/40400000-0000-0000-0000-000000000000",
			statusCode:   http.StatusNotFound,
			bodyResponse: `{"error":"foo not found"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("40400000-0000-0000-0000-000000000000"),
				).Return(
					(*model.Foo)(nil),
					repository.NewNotFound("foo", "40400000-0000-0000-0000-000000000000"),
				)
			},
		},
		{
			name:         "Failure Case - Repository Error",
			url:          "/foos/40000000-0000-0000-0000-000000000000",
			statusCode:   http.StatusInternalServerError,
			bodyResponse: `{"error":"failed to get foo by id"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("40000000-0000-0000-0000-000000000000"),
				).Return(
					(*model.Foo)(nil),
					errors.New("repository error"),
				)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			controller := NewFooController(mockHandler)

			testCase.setupMockHandler(mockHandler)

			req, err := http.NewRequest(http.MethodGet, testCase.url, nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.GET("/foos/:id", controller.GetByID)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
			assert.JSONEq(t, testCase.bodyResponse, w.Body.String())
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestFooController_Create(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		body         string
		statusCode   int
		bodyResponse string

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name:         "Success Case",
			url:          "/foos",
			body:         `{"label":"foo_create", "secret":"secret_create", "value":1, "weight":1.5}`,
			statusCode:   http.StatusCreated,
			bodyResponse: `{"id":"20000000-0000-0000-0000-000000000001"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"Create",
					mock.Anything,
					data.FooCreateInput{
						Label:  "foo_create",
						Secret: "secret_create",
						Value:  1,
						Weight: 1.5,
					}).Return(
					&model.Foo{
						Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:     "foo_create",
						Secret:    "secret_create",
						Value:     1,
						Weight:    1.5,
						CreatedAt: time.Now(),
						UpdatedAt: nil,
					}, nil)
			},
		},
		{
			name:             "Failure Case - Invalid Body",
			url:              "/foos",
			body:             `{"label":"foo_create"}`,
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate request body"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:         "Failure Case - Repository Error",
			url:          "/foos",
			body:         `{"label":"foo_create", "secret":"secret_create", "value":1, "weight":1.5}`,
			statusCode:   http.StatusInternalServerError,
			bodyResponse: `{"error":"failed to create foo"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"Create",
					mock.Anything,
					data.FooCreateInput{
						Label:  "foo_create",
						Secret: "secret_create",
						Value:  1,
						Weight: 1.5,
					}).Return(
					(*model.Foo)(nil),
					errors.New("repository error"),
				)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			controller := NewFooController(mockHandler)

			testCase.setupMockHandler(mockHandler)

			req, err := http.NewRequest(http.MethodPost, testCase.url, strings.NewReader(testCase.body))
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.POST("/foos", controller.Create)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
			assert.JSONEq(t, testCase.bodyResponse, w.Body.String())
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestFooController_Update(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		body         string
		statusCode   int
		bodyResponse string

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name:       "Success Case",
			url:        "/foos/20000000-0000-0000-0000-000000000001",
			body:       `{"label":"foo_update", "secret":"secret_update", "value":1, "weight":1.5}`,
			statusCode: http.StatusNoContent,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"Update",
					mock.Anything,
					data.FooUpdateInput{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "foo_update",
						Secret: "secret_update",
						Value:  1,
						Weight: 1.5,
					}).Return(nil)
			},
		},
		{
			name:             "Failure Case - Invalid Body",
			url:              "/foos/20000000-0000-0000-0000-000000000001",
			body:             `{"label":"foo_update"}`,
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate request body"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:         "Failure Case - Not Found",
			url:          "/foos/40400000-0000-0000-0000-000000000000",
			body:         `{"label":"foo_update", "secret":"secret_update", "value":1, "weight":1.5}`,
			statusCode:   http.StatusNotFound,
			bodyResponse: `{"error":"foo not found"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"Update",
					mock.Anything,
					data.FooUpdateInput{
						Id:     uuid.MustParse("40400000-0000-0000-0000-000000000000"),
						Label:  "foo_update",
						Secret: "secret_update",
						Value:  1,
						Weight: 1.5,
					}).Return(repository.NewNotFound("foo", "40400000-0000-0000-0000-000000000000"))
			},
		},
		{
			name:         "Failure Case - Repository Error",
			url:          "/foos/20000000-0000-0000-0000-000000000001",
			body:         `{"label":"foo_update", "secret":"secret_update", "value":1, "weight":1.5}`,
			statusCode:   http.StatusInternalServerError,
			bodyResponse: `{"error": "failed to update foo"}`,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"Update",
					mock.Anything,
					data.FooUpdateInput{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "foo_update",
						Secret: "secret_update",
						Value:  1,
						Weight: 1.5,
					}).Return(errors.New("repository error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			controller := NewFooController(mockHandler)

			testCase.setupMockHandler(mockHandler)

			req, err := http.NewRequest(http.MethodPut, testCase.url, strings.NewReader(testCase.body))
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.PUT("/foos/:id", controller.Update)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
			if testCase.bodyResponse != "" {
				assert.JSONEq(t, testCase.bodyResponse, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestFooController_Delete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		statusCode   int
		bodyResponse string

		setupMockHandler func(*service.MockFooService)
	}{
		{
			name:       "Success Case",
			url:        "/foos/20000000-0000-0000-0000-000000000001",
			statusCode: http.StatusNoContent,

			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
		},
		{
			name:             "Failure Case - Not UUID",
			url:              "/foos/not_uuid",
			statusCode:       http.StatusBadRequest,
			bodyResponse:     `{"error":"failed to validate path params"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {},
		},
		{
			name:         "Failure Case - Not Found",
			url:          "/foos/40400000-0000-0000-0000-000000000000",
			statusCode:   http.StatusNotFound,
			bodyResponse: `{"error":"foo not found"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("40400000-0000-0000-0000-000000000000"),
				).Return(repository.NewNotFound("foo", "40400000-0000-0000-0000-000000000000"))
			},
		},
		{
			name:         "Failure Case - Repository Error",
			url:          "/foos/20000000-0000-0000-0000-000000000001",
			statusCode:   http.StatusInternalServerError,
			bodyResponse: `{"error": "failed to delete foo"}`,
			setupMockHandler: func(mockHandler *service.MockFooService) {
				mockHandler.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(errors.New("repository error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(service.MockFooService)
			controller := NewFooController(mockHandler)

			testCase.setupMockHandler(mockHandler)

			req, err := http.NewRequest(http.MethodDelete, testCase.url, nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.DELETE("/foos/:id", controller.DeleteByID)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Code)
			if testCase.bodyResponse != "" {
				assert.JSONEq(t, testCase.bodyResponse, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
			mockHandler.AssertExpectations(t)
		})
	}
}
