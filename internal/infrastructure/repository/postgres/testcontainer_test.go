package postgres

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
	Config Config
}

// CreatePostgresContainer initializes and starts a PostgreSQL container with predefined configurations and returns its instance.
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	config := Config{
		Host:          "localhost",
		User:          "postgres",
		Password:      "postgres",
		DBName:        "test-db",
		SSLMode:       "disable",
		MaxOpenConns:  10,
		MaxIdleConns:  5,
		MaxLifetime:   10,
		Migrate:       true,
		MigrationPath: "file://../../../../migrations/postgres",
	}
	pgContainer, err := postgres.Run(ctx, "postgres:15.3-alpine",
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

func seed(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier seed.sql: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erreur lors du début de la transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("erreur lors de l'exécution du seed: %w", err)
	}

	return tx.Commit()
}
