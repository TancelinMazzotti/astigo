package nats

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

// NatsContainer is a type that embeds a NATS container and its configuration for managing NATS server instances.
type NatsContainer struct {
	*nats.NATSContainer
	Config NatsConfig
}

// CreateNatsContainer creates and initializes a NATS container in a Docker environment with the provided context.
// Returns a pointer to a NatsContainer instance containing the container instance and its configuration, or an error.
func CreateNatsContainer(ctx context.Context) (*NatsContainer, error) {
	config := NatsConfig{
		URL:      "localhost",
		Username: "foo",
		Password: "bar",
	}

	natsContainer, err := nats.Run(ctx, "nats:2.9",
		nats.WithUsername(config.Username),
		nats.WithPassword(config.Password),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("4222/tcp").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	uri, err := natsContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	config.URL = uri

	return &NatsContainer{
		NATSContainer: natsContainer,
		Config:        config,
	}, nil
}
