package event

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ConsumerNats struct {
	Logger *zap.Logger
	conn   *nats.Conn

	fooWorker *FooWorkerNats
}

func (c *ConsumerNats) Close() error {
	if err := c.fooWorker.Close(); err != nil {
		return fmt.Errorf("failed to close foo consumer: %w", err)
	}
	return nil
}

func NewConsumerNats(logger *zap.Logger, conn *nats.Conn) (*ConsumerNats, error) {
	var err error
	consumer := &ConsumerNats{
		Logger: logger,
		conn:   conn,
	}

	if consumer.fooWorker, err = NewFooWorkerNats(logger, conn, "foo"); err != nil {
		return nil, fmt.Errorf("fail to create foo consumer: %w", err)
	}

	return consumer, nil
}
