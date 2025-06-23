package nats

import (
	"astigo/internal/domain/messaging"
	"astigo/internal/domain/model"
	"context"
	"encoding/json"
	"fmt"
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

type FooNats struct {
	conn *nats.Conn
}

func (n *FooNats) PublishFooCreated(_ context.Context, foo model.Foo) error {
	data, err := json.Marshal(foo)
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation du Foo: %w", err)
	}

	return n.conn.Publish(fooCreatedSubject, data)
}

func (n *FooNats) PublishFooUpdated(_ context.Context, foo model.Foo) error {
	data, err := json.Marshal(foo)
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation du Foo: %w", err)
	}

	return n.conn.Publish(fooUpdatedSubject, data)
}

func (n *FooNats) PublishFooDeleted(_ context.Context, id uuid.UUID) error {
	data, err := json.Marshal(map[string]string{"id": id.String()})
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation de l'ID: %w", err)
	}

	return n.conn.Publish(fooDeletedSubject, data)
}

func NewFooNats(conn *nats.Conn) *FooNats {
	return &FooNats{conn: conn}
}
