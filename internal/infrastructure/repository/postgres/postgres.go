package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel/attribute"
)

// Config represents the configuration settings required to connect to a PostgreSQL database.
type Config struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	DBName        string `mapstructure:"dbname"`
	SSLMode       string `mapstructure:"sslmode"`
	MaxOpenConns  int    `mapstructure:"max_open_conns"`
	MaxIdleConns  int    `mapstructure:"max_idle_conns"`
	MaxLifetime   int    `mapstructure:"max_lifetime"`
	Migrate       bool   `mapstructure:"migrate"`
	MigrationPath string `mapstructure:"migration_path"`
}

// NewPostgres initializes and returns a PostgreSQL database connection based on the provided configuration.
// It configures connection pool settings and verifies connectivity by pinging the database.
func NewPostgres(ctx context.Context, config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode,
	)

	db, err := otelsql.Open("pgx", dsn,
		otelsql.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.host", config.Host),
			attribute.Int("db.port", config.Port),
			attribute.String("db.name", config.DBName),
			attribute.String("db.user", config.User),
		),
		otelsql.WithDBName(config.DBName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	if config.Migrate {
		if err := migration(db, config.MigrationPath); err != nil {
			return nil, fmt.Errorf("failed to migrate postgres: %w", err)
		}
	}

	return db, nil
}

func migration(db *sql.DB, source string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		source,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create new migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to migrate postgres: %w", err)
	}

	return nil
}
