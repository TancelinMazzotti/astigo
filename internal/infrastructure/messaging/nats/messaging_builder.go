package nats

import (
	"context"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/messaging"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/connector"
	"github.com/nats-io/nats.go"
)

const Driver = "nats"

var (
	_ messaging.IMessagingBuilder = (*MessagingBuilder)(nil)
)

type MessagingBuilder struct {
	config connector.NatsConfig
	conn   *nats.Conn
	foo    *FooNats
}

func (m *MessagingBuilder) Build(ctx context.Context) error {
	var err error

	m.conn, err = connector.NewNats(m.config)
	if err != nil {
		return fmt.Errorf("failed to create nats client: %w", err)
	}

	m.foo = NewFooNats(m.conn)

	return nil
}

func (m *MessagingBuilder) GetFooMessaging() (messaging.IFooMessaging, error) {
	if m.foo == nil {
		return nil, fmt.Errorf("foo messaging is not initialized")
	}
	return m.foo, nil
}

func NewMessagingBuilder(config connector.NatsConfig) messaging.IMessagingBuilder {
	return &MessagingBuilder{config: config}
}
