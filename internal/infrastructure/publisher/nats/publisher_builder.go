package nats

import (
	"context"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/publisher"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/connector"
	"github.com/nats-io/nats.go"
)

const Driver = "nats"

var (
	_ publisher.IPublisherBuilder = (*PublisherBuilder)(nil)
)

type PublisherBuilder struct {
	config connector.NatsConfig
	conn   *nats.Conn
	foo    *FooNats
}

func (m *PublisherBuilder) Build(ctx context.Context) error {
	var err error

	m.conn, err = connector.NewNats(m.config)
	if err != nil {
		return fmt.Errorf("failed to create nats client: %w", err)
	}

	m.foo = NewFooNats(m.conn)

	return nil
}

func (m *PublisherBuilder) GetFooPublisher() (publisher.IFooPublisher, error) {
	if m.foo == nil {
		return nil, fmt.Errorf("foo publisher is not initialized")
	}
	return m.foo, nil
}

func NewPublisherBuilder(config connector.NatsConfig) publisher.IPublisherBuilder {
	return &PublisherBuilder{config: config}
}
