package nats

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/messaging/nats/message"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationFooNats_PublishFooCreated tests the FooNats integration by publishing a "foo.created" message to NATS.
func TestIntegrationFooNats_PublishFooCreated(t *testing.T) {
	t.Parallel()
	now := time.Now()
	testCases := []struct {
		name          string
		foo           *model.Foo
		expectedError error
		receivedData  message.FooMessage
	}{
		{
			name: "Success Case",
			foo: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "foo_create",
				Secret:    "secret_create",
				Value:     10,
				Weight:    1.5,
				CreatedAt: now,
				UpdatedAt: nil,
				Bars:      nil,
			},
			expectedError: nil,
			receivedData: message.FooMessage{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "foo_create",
				Value:     10,
				Weight:    1.5,
				CreatedAt: now,
				UpdatedAt: nil,
			},
		},
	}
	ctx := context.Background()
	container, err := CreateNatsContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	nc, err := NewNats(container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			messageChan := make(chan message.FooMessage, 1)
			sub, err := nc.Subscribe(fooCreatedSubject, func(msg *nats.Msg) {
				var receivedDta message.FooMessage
				err := json.Unmarshal(msg.Data, &receivedDta)
				if err != nil {
					t.Error("failed to unmarshal data:", err)
					return
				}
				messageChan <- receivedDta
			})
			if err != nil {
				t.Fatal("failed to subscribe:", err)
			}
			defer sub.Unsubscribe()

			messaging := NewFooNats(nc)

			err = messaging.PublishFooCreated(ctx, testCase.foo)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			select {
			case receivedFoo := <-messageChan:
				assert.Equal(t, testCase.foo.Id, receivedFoo.Id)
				assert.Equal(t, testCase.foo.Label, receivedFoo.Label)
				assert.Equal(t, testCase.foo.Value, receivedFoo.Value)
				assert.Equal(t, testCase.foo.Weight, receivedFoo.Weight)
				// We allow a one-second difference because NATS internal conversions
				// (serialization/deserialization) may slightly modify the timestamp precision
				assert.WithinDuration(t, testCase.foo.CreatedAt, receivedFoo.CreatedAt, time.Second)
				assert.Nil(t, receivedFoo.UpdatedAt)

			case <-time.After(2 * time.Second):
				t.Error("timeout: no message received")
			}

		})
	}
}

// TestIntegrationFooNats_PublishFooUpdated tests the PublishFooUpdated method by verifying NATS message delivery and content.
func TestIntegrationFooNats_PublishFooUpdated(t *testing.T) {
	t.Parallel()
	now := time.Now()
	created := now.Add(-1 * time.Hour)
	testCases := []struct {
		name          string
		foo           *model.Foo
		expectedError error
		receivedData  message.FooMessage
	}{
		{
			name: "Success Case",
			foo: &model.Foo{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "foo_updated",
				Secret:    "secret_updated",
				Value:     10,
				Weight:    1.5,
				CreatedAt: created,
				UpdatedAt: &now,
			},
			expectedError: nil,
			receivedData: message.FooMessage{
				Id:        uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:     "foo_updated",
				Value:     10,
				Weight:    1.5,
				CreatedAt: created,
				UpdatedAt: &now,
			},
		},
	}
	ctx := context.Background()
	container, err := CreateNatsContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	nc, err := NewNats(container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			messageChan := make(chan message.FooMessage, 1)
			sub, err := nc.Subscribe(fooUpdatedSubject, func(msg *nats.Msg) {
				var receivedFoo message.FooMessage
				err := json.Unmarshal(msg.Data, &receivedFoo)
				if err != nil {
					t.Error("failed to unmarshal data:", err)
					return
				}
				messageChan <- receivedFoo
			})
			if err != nil {
				t.Fatal("failed to subscribe:", err)
			}
			defer sub.Unsubscribe()

			messaging := NewFooNats(nc)

			err = messaging.PublishFooUpdated(ctx, testCase.foo)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			select {
			case receivedFoo := <-messageChan:
				assert.Equal(t, testCase.foo.Id, receivedFoo.Id)
				assert.Equal(t, testCase.foo.Label, receivedFoo.Label)
				assert.Equal(t, testCase.foo.Value, receivedFoo.Value)
				assert.Equal(t, testCase.foo.Weight, receivedFoo.Weight)
				// We allow a one-second difference because NATS internal conversions
				// (serialization/deserialization) may slightly modify the timestamp precision
				assert.WithinDuration(t, testCase.foo.CreatedAt, receivedFoo.CreatedAt, time.Second)
				assert.WithinDuration(t, *testCase.foo.UpdatedAt, *receivedFoo.UpdatedAt, time.Second)
			case <-time.After(2 * time.Second):
				t.Error("timeout: no message received")
			}

		})
	}
}

// TestIntegrationFooNats_PublishFooDeleted validates that the "foo.deleted" message is published correctly to the NATS server.
func TestIntegrationFooNats_PublishFooDeleted(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		id            uuid.UUID
		expectedError error
		receivedID    uuid.UUID
	}{
		{
			name:          "Success Case",
			id:            uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			expectedError: nil,
			receivedID:    uuid.MustParse("20000000-0000-0000-0000-000000000001"),
		},
	}
	ctx := context.Background()
	container, err := CreateNatsContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	nc, err := NewNats(container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			messageChan := make(chan uuid.UUID, 1)
			sub, err := nc.Subscribe(fooDeletedSubject, func(msg *nats.Msg) {
				var receivedId struct {
					Id uuid.UUID
				}
				err := json.Unmarshal(msg.Data, &receivedId)
				if err != nil {
					t.Error("failed to unmarshal data:", err)
					return
				}
				messageChan <- receivedId.Id
			})
			if err != nil {
				t.Fatal("failed to subscribe:", err)
			}
			defer sub.Unsubscribe()

			messaging := NewFooNats(nc)

			err = messaging.PublishFooDeleted(ctx, testCase.id)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			select {
			case receivedID := <-messageChan:
				assert.Equal(t, testCase.receivedID, receivedID)
			case <-time.After(2 * time.Second):
				t.Error("timeout: no message received")
			}

		})
	}
}
