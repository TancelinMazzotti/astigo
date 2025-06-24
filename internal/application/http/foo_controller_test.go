package http

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFooController_GetAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		statusCode   int
		bodyResponse string

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name:       "Success Case - Multiple Foos",
			url:        "/foos?offset=0&limit=10",
			statusCode: http.StatusOK,
			bodyResponse: `[
				{"id":"20000000-0000-0000-0000-000000000001", "label":"Foo1"},
				{"id":"20000000-0000-0000-0000-000000000002", "label":"Foo2"},
				{"id":"20000000-0000-0000-0000-000000000003", "label":"Foo3"}
			]`,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
				mockHandler.On(
					"GetAll",
					mock.Anything,
					handler.FooReadListInput{Offset: 0, Limit: 10},
				).Return([]model.Foo{
					{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "Foo1"},
					{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "Foo2"},
					{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "Foo3"},
				}, nil)
			},
		},
		{
			name:         "Failure Case - Repository Error",
			url:          "/foos?offset=0&limit=10",
			statusCode:   http.StatusInternalServerError,
			bodyResponse: `{"error":"repository error"}`,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
				mockHandler.On(
					"GetAll",
					mock.Anything,
					handler.FooReadListInput{Offset: 0, Limit: 10},
				).Return(([]model.Foo)(nil), errors.New("repository error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
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

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name:         "Success Case",
			url:          "/foos/20000000-0000-0000-0000-000000000001",
			statusCode:   http.StatusOK,
			bodyResponse: `{"id":"20000000-0000-0000-0000-000000000001", "label":"Foo1"}`,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
				mockHandler.On(
					"GetByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(
					&model.Foo{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "Foo1",
						Secret: "secret1",
					}, nil)
			},
		},
		{
			name:         "Failure Case - Not Found",
			url:          "/foos/40400000-0000-0000-0000-000000000000",
			statusCode:   http.StatusNotFound,
			bodyResponse: `{"error":"foo with id '40400000-0000-0000-0000-000000000000' not found"}`,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
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
			bodyResponse: `{"error":"repository error"}`,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
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
			mockHandler := new(handler.MockFooHandler)
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

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name:         "Success Case",
			url:          "/foos",
			body:         `{"label":"foo_create", "secret":"secret_create"}`,
			statusCode:   http.StatusCreated,
			bodyResponse: `{"id":"20000000-0000-0000-0000-000000000001"}`,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
				mockHandler.On(
					"Create",
					mock.Anything,
					handler.FooCreateInput{
						Label:  "foo_create",
						Secret: "secret_create",
					}).Return(
					&model.Foo{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "foo_create",
						Secret: "secret_create",
					}, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
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
		name       string
		url        string
		body       string
		statusCode int

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name:       "Success Case",
			url:        "/foos/20000000-0000-0000-0000-000000000001",
			body:       `{"label":"foo_update", "secret":"secret_update"}`,
			statusCode: http.StatusNoContent,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
				mockHandler.On(
					"Update",
					mock.Anything,
					handler.FooUpdateInput{
						Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
						Label:  "foo_update",
						Secret: "secret_update",
					}).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
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
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestFooController_Delete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		url        string
		statusCode int

		setupMockHandler func(*handler.MockFooHandler)
	}{
		{
			name:       "Success Case",
			url:        "/foos/20000000-0000-0000-0000-000000000001",
			statusCode: http.StatusNoContent,

			setupMockHandler: func(mockHandler *handler.MockFooHandler) {
				mockHandler.On(
					"DeleteByID",
					mock.Anything,
					uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockHandler := new(handler.MockFooHandler)
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
			mockHandler.AssertExpectations(t)

		})
	}
}
