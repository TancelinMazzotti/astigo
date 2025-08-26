package event

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	fooCreatedSubject = "foo.created"
	fooUpdatedSubject = "foo.updated"
	fooDeletedSubject = "foo.deleted"
)

type FooWorkerNats struct {
	Logger        *zap.Logger
	conn          *nats.Conn
	subscriptions []*nats.Subscription
}

func (f *FooWorkerNats) OnCreated(msg *nats.Msg) {
	f.Logger.Info("on created", zap.String("msg", string(msg.Data)))
}

func (f *FooWorkerNats) OnUpdated(msg *nats.Msg) {
	f.Logger.Info("on updated", zap.String("msg", string(msg.Data)))
}

func (f *FooWorkerNats) OnDeleted(msg *nats.Msg) {
	f.Logger.Info("on deleted", zap.String("msg", string(msg.Data)))
}

func (f *FooWorkerNats) Close() error {
	for _, sub := range f.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

func NewFooWorkerNats(logger *zap.Logger, conn *nats.Conn, group string) (*FooWorkerNats, error) {
	foo := &FooWorkerNats{
		Logger:        logger,
		conn:          conn,
		subscriptions: []*nats.Subscription{},
	}

	sub, err := conn.QueueSubscribe(fooCreatedSubject, group, foo.OnCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to foo.created: %w", err)
	}
	foo.subscriptions = append(foo.subscriptions, sub)

	sub, err = conn.QueueSubscribe(fooUpdatedSubject, group, foo.OnUpdated)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to foo.updated: %w", err)
	}
	foo.subscriptions = append(foo.subscriptions, sub)

	sub, err = conn.QueueSubscribe(fooDeletedSubject, group, foo.OnDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to foo.deleted: %w", err)
	}
	foo.subscriptions = append(foo.subscriptions, sub)

	return foo, nil
}
