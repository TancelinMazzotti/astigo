package http

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
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
		name          string
		queryParams   string
		mockRequest   handler.PaginationInput
		mockResponse  []model.Foo
		mockError     error
		statusCode    int
		bodyResponse  string
		expectedError error
	}{
		{
			name:        "Success Case - Multiple Foos",
			queryParams: "offset=0&limit=10",
			mockRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockResponse: []model.Foo{
				{Id: 1, Label: "Foo1"},
				{Id: 2, Label: "Foo2"},
				{Id: 3, Label: "Foo3"},
			},
			mockError:  nil,
			statusCode: http.StatusOK,
			bodyResponse: `[
				{"id":1, "label":"Foo1"},
				{"id":2, "label":"Foo2"},
				{"id":3, "label":"Foo3"}
			]`,
			expectedError: nil,
		},
		{
			name:        "Failure Case - Repository Error",
			queryParams: "offset=0&limit=10",
			mockRequest: handler.PaginationInput{
				Offset: 0,
				Limit:  10,
			},
			mockResponse:  nil,
			mockError:     errors.New("repository error"),
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

			req, err := http.NewRequest(http.MethodGet, "/foos?"+testCase.queryParams, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			controller.GetAll(c)

			assert.Equal(t, testCase.statusCode, w.Code)
			assert.JSONEq(t, testCase.bodyResponse, w.Body.String())
			mockHandler.AssertExpectations(t)
		})
	}
}
