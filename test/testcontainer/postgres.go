package testcontainer

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer represents a PostgreSQL container with its configuration settings encapsulated.
type PostgresContainer struct {
	*postgres.PostgresContainer
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// CreatePostgresContainer initializes and starts a PostgreSQL container with predefined configurations and returns its instance.
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	user := "postgres"
	password := "postgres"
	database := "test-db"
	pgContainer, err := postgres.Run(ctx, "postgres:15.3-alpine",
		postgres.WithDatabase(database),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
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

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		Host:              "localhost",
		Port:              containerPort.Int(),
		User:              user,
		Password:          password,
		Database:          database,
	}, nil
}

func SeedPostgres(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("fail to read seed.sql: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("fail to start transcation: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("fail to execute seed: %w", err)
	}

	return tx.Commit()
}
