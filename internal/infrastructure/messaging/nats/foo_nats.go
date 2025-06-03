package nats

import (
	"astigo/internal/domain/messaging"
	"github.com/nats-io/nats.go"
)

var (
	_ messaging.IFooMessaging = (*FooNats)(nil)
)

type FooNats struct {
	conn *nats.Conn
}

func NewFooNats(conn *nats.Conn) *FooNats {
	return &FooNats{conn: conn}
}
