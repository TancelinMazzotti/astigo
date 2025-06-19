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
	"testing"
)

func TestFooController_GetAll(t *testing.T) {
	testCases := []struct {
		name string
		url  string

		mockRequest  handler.PaginationInput
		mockResponse []model.Foo
		mockError    error

		statusCode    int
		bodyResponse  string
		expectedError error
	}{
		{
			name: "Success Case - Multiple Foos",
			url:  "/foos?offset=0&limit=10",

			mockRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockResponse: []model.Foo{
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000001"), Label: "Foo1"},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000002"), Label: "Foo2"},
				{Id: uuid.MustParse("20000000-0000-0000-0000-000000000003"), Label: "Foo3"},
			},
			mockError: nil,

			statusCode: http.StatusOK,
			bodyResponse: `[
				{"id":"20000000-0000-0000-0000-000000000001", "label":"Foo1"},
				{"id":"20000000-0000-0000-0000-000000000002", "label":"Foo2"},
				{"id":"20000000-0000-0000-0000-000000000003", "label":"Foo3"}
			]`,
			expectedError: nil,
		},
		{
			name: "Failure Case - Repository Error",
			url:  "/foos?offset=0&limit=10",

			mockRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockResponse: nil,
			mockError:    errors.New("repository error"),

			statusCode:    http.StatusInternalServerError,
			bodyResponse:  `{"error":"repository error"}`,
			expectedError: errors.New("fail to find all foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("GetAll", mock.Anything, testCase.mockRequest).Return(testCase.mockResponse, testCase.mockError)
			controller := NewFooController(mockHandler)

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
	testCases := []struct {
		name string
		url  string

		mockRequest  uuid.UUID
		mockResponse *model.Foo
		mockError    error

		statusCode    int
		bodyResponse  string
		expectedError error
	}{
		{
			name: "Success Case",
			url:  "/foos/20000000-0000-0000-0000-000000000001",

			mockRequest: uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			mockResponse: &model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "Foo1",
				Secret: "secret1",
			},
			mockError: nil,

			statusCode:    http.StatusOK,
			bodyResponse:  `{"id":"20000000-0000-0000-0000-000000000001", "label":"Foo1"}`,
			expectedError: nil,
		},
		{
			name: "Failure Case - Not Found",
			url:  "/foos/40400000-0000-0000-0000-000000000000",

			mockRequest:  uuid.MustParse("40400000-0000-0000-0000-000000000000"),
			mockResponse: nil,
			mockError:    repository.NewNotFound("foo", "40400000-0000-0000-0000-000000000000"),

			statusCode:    http.StatusNotFound,
			bodyResponse:  `{"error":"foo with id '40400000-0000-0000-0000-000000000000' not found"}`,
			expectedError: errors.New("fail to find foo by id: repository error"),
		},
		{
			name: "Failure Case - Repository Error",
			url:  "/foos/40000000-0000-0000-0000-000000000000",

			mockRequest:  uuid.MustParse("40000000-0000-0000-0000-000000000000"),
			mockResponse: nil,
			mockError:    errors.New("repository error"),

			statusCode:    http.StatusInternalServerError,
			bodyResponse:  `{"error":"repository error"}`,
			expectedError: errors.New("fail to find foo by id: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("GetByID", mock.Anything, testCase.mockRequest).Return(testCase.mockResponse, testCase.mockError)
			controller := NewFooController(mockHandler)

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
