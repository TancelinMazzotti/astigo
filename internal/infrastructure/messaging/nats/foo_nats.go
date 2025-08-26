package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/contract/messaging"
	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/messaging/nats/message"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	fooCreatedSubject = "foo.created"
	fooUpdatedSubject = "foo.updated"
	fooDeletedSubject = "foo.deleted"
)

var (
	_ messaging.IFooMessaging = (*FooNats)(nil)
)

// FooNats wraps a NATS connection and implements the IFooMessaging interface for publishing Foo-related messages.
type FooNats struct {
	conn *nats.Conn
}

// PublishFooCreated publishes a "foo.created" message to the NATS server using the provided Foo data.
func (n *FooNats) PublishFooCreated(_ context.Context, foo *model.Foo) error {
	msg := message.NewFooMessage(foo)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize Foo: %w", err)
	}

	return n.conn.Publish(fooCreatedSubject, data)
}

// PublishFooUpdated publishes a "foo.updated" message to the NATS server with the updated Foo data.
func (n *FooNats) PublishFooUpdated(_ context.Context, foo *model.Foo) error {
	msg := message.NewFooMessage(foo)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize Foo: %w", err)
	}

	return n.conn.Publish(fooUpdatedSubject, data)
}

// PublishFooDeleted publishes a "foo.deleted" message to the NATS server with the provided Foo ID serialized as JSON.
func (n *FooNats) PublishFooDeleted(_ context.Context, id uuid.UUID) error {
	data, err := json.Marshal(map[string]string{"id": id.String()})
	if err != nil {
		return fmt.Errorf("failed to serialize ID: %w", err)
	}

	return n.conn.Publish(fooDeletedSubject, data)
}

func NewFooNats(conn *nats.Conn) *FooNats {
	return &FooNats{conn: conn}
}
