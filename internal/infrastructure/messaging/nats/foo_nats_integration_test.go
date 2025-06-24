package nats

import (
	"astigo/internal/domain/model"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntegrationFooNats_PublishFooCreated(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		foo           model.Foo
		expectedError error
		receivedFoo   model.Foo
	}{
		{
			name: "Success Case",
			foo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_create",
				Secret: "secret_create",
			},
			expectedError: nil,
			receivedFoo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_create",
				Secret: "secret_create",
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
			messageChan := make(chan model.Foo, 1)
			sub, err := nc.Subscribe(fooCreatedSubject, func(msg *nats.Msg) {
				var receivedFoo model.Foo
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

			err = messaging.PublishFooCreated(ctx, testCase.foo)

			if testCase.expectedError != nil {
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			select {
			case receivedFoo := <-messageChan:
				assert.Equal(t, testCase.receivedFoo.Id, receivedFoo.Id)
				assert.Equal(t, testCase.receivedFoo.Label, receivedFoo.Label)
				assert.Equal(t, testCase.receivedFoo.Secret, receivedFoo.Secret)
			case <-time.After(2 * time.Second):
				t.Error("timeout: no message received")
			}

		})
	}
}

func TestIntegrationFooNats_PublishFooUpdated(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		foo           model.Foo
		expectedError error
		receivedFoo   model.Foo
	}{
		{
			name: "Success Case",
			foo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_updated",
				Secret: "secret_updated",
			},
			expectedError: nil,
			receivedFoo: model.Foo{
				Id:     uuid.MustParse("20000000-0000-0000-0000-000000000001"),
				Label:  "foo_updated",
				Secret: "secret_updated",
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
			messageChan := make(chan model.Foo, 1)
			sub, err := nc.Subscribe(fooUpdatedSubject, func(msg *nats.Msg) {
				var receivedFoo model.Foo
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
				assert.Equal(t, testCase.receivedFoo.Id, receivedFoo.Id)
				assert.Equal(t, testCase.receivedFoo.Label, receivedFoo.Label)
				assert.Equal(t, testCase.receivedFoo.Secret, receivedFoo.Secret)
			case <-time.After(2 * time.Second):
				t.Error("timeout: no message received")
			}

		})
	}
}

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
