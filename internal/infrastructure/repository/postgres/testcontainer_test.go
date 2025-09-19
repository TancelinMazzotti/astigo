package postgres

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer represents a PostgreSQL container with its configuration settings encapsulated.
type PostgresContainer struct {
	*postgres.PostgresContainer
	Config Config
}

// CreatePostgresContainer initializes and starts a PostgreSQL container with predefined configurations and returns its instance.
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	config := Config{
		Host:         "localhost",
		User:         "postgres",
		Password:     "postgres",
		DBName:       "test-db",
		SSLMode:      "disable",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		MaxLifetime:  10,
	}
	pgContainer, err := postgres.Run(ctx, "postgres:15.3-alpine",
		postgres.WithInitScripts("init-db.sql"),
		postgres.WithDatabase(config.DBName),
		postgres.WithUsername(config.User),
		postgres.WithPassword(config.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	containerPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}
	config.Port = containerPort.Int()

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		Config:            config,
	}, nil
}
