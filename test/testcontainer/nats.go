package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NatsContainer is a type that embeds a NATS container and its configuration for managing NATS server instances.
type NatsContainer struct {
	*nats.NATSContainer
	URL      string
	Username string
	Password string
}

// CreateNatsContainer creates and initializes a NATS container in a Docker environment with the provided context.
// Returns a pointer to a NatsContainer instance containing the container instance and its configuration, or an error.
func CreateNatsContainer(ctx context.Context) (*NatsContainer, error) {
	username := "foo"
	password := "bar"

	natsContainer, err := nats.Run(ctx, "nats:2.9",
		nats.WithUsername(username),
		nats.WithPassword(password),
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

	return &NatsContainer{
		NATSContainer: natsContainer,
		URL:           uri,
		Username:      username,
		Password:      password,
	}, nil
}
