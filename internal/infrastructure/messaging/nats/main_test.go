package nats

import (
	"context"
	"testing"

	"github.com/TancelinMazzotti/astigo/internal/infrastructure/connector"
	"github.com/TancelinMazzotti/astigo/test/testcontainer"
	"github.com/nats-io/nats.go"
)

var (
	DefaultConfig        connector.NatsConfig
	DefaultReadContainer *testcontainer.NatsContainer

	DefaultConnector *nats.Conn
)

func TestMain(m *testing.M) {
	var err error

	ctx := context.Background()
	DefaultReadContainer, err = testcontainer.CreateNatsContainer(ctx)
	if err != nil {
		panic(err)
	}

	DefaultConfig = connector.NatsConfig{
		URL:      DefaultReadContainer.URL,
		Username: DefaultReadContainer.User,
		Password: DefaultReadContainer.Password,
	}

	DefaultConnector, err = connector.NewNats(DefaultConfig)
	if err != nil {
		panic(err)
	}

	m.Run()
}
