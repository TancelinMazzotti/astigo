package http

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"errors"
	"github.com/gin-gonic/gin"
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
				{Id: 1, Label: "Foo1"},
				{Id: 2, Label: "Foo2"},
				{Id: 3, Label: "Foo3"},
			},
			mockError: nil,

			statusCode: http.StatusOK,
			bodyResponse: `[
				{"id":1, "label":"Foo1"},
				{"id":2, "label":"Foo2"},
				{"id":3, "label":"Foo3"}
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

		mockRequest  int
		mockResponse *model.Foo
		mockError    error

		statusCode    int
		bodyResponse  string
		expectedError error
	}{
		{
			name: "Success Case",
			url:  "/foos/1",

			mockRequest: 1,
			mockResponse: &model.Foo{
				Id:     1,
				Label:  "Foo1",
				Secret: "secret1",
			},
			mockError: nil,

			statusCode:    http.StatusOK,
			bodyResponse:  `{"id":1, "label":"Foo1"}`,
			expectedError: nil,
		},
		{
			name: "Failure Case - Not Found",
			url:  "/foos/-1",

			mockRequest:  -1,
			mockResponse: nil,
			mockError:    repository.NewNotFound("foo", "-1"),

			statusCode:    http.StatusNotFound,
			bodyResponse:  `{"error":"foo with id '-1' not found"}`,
			expectedError: errors.New("fail to find foo by id: repository error"),
		},
		{
			name: "Failure Case - Repository Error",
			url:  "/foos/-1",

			mockRequest:  -1,
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
