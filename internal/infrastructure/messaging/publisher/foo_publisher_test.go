package publisher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/test/mocks/domain/port/out/mockmessaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFooPublisher_PublishFooCreated(t *testing.T) {
	t.Parallel()
	createdAt := time.Now()

	testCases := []struct {
		name          string
		foo           *model.Foo
		expectedError error

		setupMockMessaging1 func(*mockmessaging.MockFooMessaging)
		setupMockMessaging2 func(*mockmessaging.MockFooMessaging)
	}{
		{
			name: "Success Case",
			foo: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "Foo1",
				Secret:    "secret1",
				Value:     1,
				Weight:    1.5,
				CreatedAt: createdAt,
			},
			expectedError: nil,

			setupMockMessaging1: func(mockMessaging *mockmessaging.MockFooMessaging) {
				mockMessaging.On("PublishFooCreated", mock.Anything, &model.Foo{
					Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:     "Foo1",
					Secret:    "secret1",
					Value:     1,
					Weight:    1.5,
					CreatedAt: createdAt,
				}).Return(nil)
			},
			setupMockMessaging2: func(mockMessaging *mockmessaging.MockFooMessaging) {
				mockMessaging.On("PublishFooCreated", mock.Anything, &model.Foo{
					Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:     "Foo1",
					Secret:    "secret1",
					Value:     1,
					Weight:    1.5,
					CreatedAt: createdAt,
				}).Return(nil)
			},
		},
		{
			name: "Failure Case - All Messaging Error",
			foo: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "Foo1",
				Secret:    "secret1",
				Value:     1,
				Weight:    1.5,
				CreatedAt: createdAt,
			},
			expectedError: errors.New("failed to publish foo created: failed to publish"),

			setupMockMessaging1: func(mockMessaging *mockmessaging.MockFooMessaging) {
				mockMessaging.On("PublishFooCreated", mock.Anything, &model.Foo{
					Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:     "Foo1",
					Secret:    "secret1",
					Value:     1,
					Weight:    1.5,
					CreatedAt: createdAt,
				}).Return(errors.New("failed to publish"))
			},
			setupMockMessaging2: func(mockMessaging *mockmessaging.MockFooMessaging) {
				mockMessaging.On("PublishFooCreated", mock.Anything, &model.Foo{
					Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
					Label:     "Foo1",
					Secret:    "secret1",
					Value:     1,
					Weight:    1.5,
					CreatedAt: createdAt,
				}).Return(errors.New("failed to publish"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockMessaging1 := new(mockmessaging.MockFooMessaging)
			mockMessaging2 := new(mockmessaging.MockFooMessaging)

			testCase.setupMockMessaging1(mockMessaging1)
			testCase.setupMockMessaging2(mockMessaging2)

			pub := NewFooPublisher()
			pub.Subscribe(mockMessaging1)
			pub.Subscribe(mockMessaging2)

			err := pub.PublishFooCreated(context.Background(), testCase.foo)

			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
