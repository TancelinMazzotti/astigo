package http

import (
	"astigo/internal/domain/handler"
	"astigo/pkg/dto"
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
		mockResponse  []dto.FooResponseReadDto
		mockError     error
		queryParams   string
		statusCode    int
		bodyResponse  string
		expectedError error
	}{
		{
			name: "Success Case - Multiple Foos",
			mockResponse: []dto.FooResponseReadDto{
				{Id: 1, Label: "Foo1", Bars: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
				{Id: 2, Label: "Foo2", Bars: []int{}},
				{Id: 3, Label: "Foo3", Bars: []int{}},
			},
			mockError:   nil,
			queryParams: "offset=0&limit=10",
			statusCode:  http.StatusOK,
			bodyResponse: `[
				{"id":1, "label":"Foo1", "bars": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]},
				{"id":2, "label":"Foo2", "bars": []},
				{"id":3, "label":"Foo3", "bars": []}
			]`,
			expectedError: nil,
		},
		{
			name:          "Failure Case - Repository Error",
			mockResponse:  nil,
			mockError:     errors.New("repository error"),
			queryParams:   "offset=0&limit=10",
			statusCode:    http.StatusInternalServerError,
			bodyResponse:  `{"error":"repository error"}`,
			expectedError: errors.New("fail to find all foo: repository error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockHandler := new(handler.MockFooHandler)
			mockHandler.On("GetAll", mock.Anything, mock.Anything).Return(testCase.mockResponse, testCase.mockError)

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
