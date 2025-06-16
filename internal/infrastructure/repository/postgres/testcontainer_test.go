package postgres

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	Config PostgresConfig
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	config := PostgresConfig{
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
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
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
