package postgres

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/TancelinMazzotti/astigo/internal/adapter/out/connector"
	"github.com/TancelinMazzotti/astigo/test/helpers"
	"github.com/TancelinMazzotti/astigo/test/testcontainer"
)

var (
	MigrationPath   string
	DefaultSeedPath string

	DefaultConfig        connector.PostgresConfig
	DefaultReadContainer *testcontainer.PostgresContainer

	DefaultReadDb *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	MigrationPath = strings.Join([]string{"file://", helpers.ProjectRoot, "/migrations/postgres"}, "")
	DefaultSeedPath = strings.Join([]string{helpers.ProjectRoot, "/test/testdata/postgres/seed.sql"}, "")

	ctx := context.Background()
	DefaultReadContainer, err = testcontainer.CreatePostgresContainer(ctx)
	if err != nil {
		panic(err)
	}

	DefaultConfig = connector.PostgresConfig{
		Host:          DefaultReadContainer.Host,
		Port:          DefaultReadContainer.Port,
		DBName:        DefaultReadContainer.Database,
		User:          DefaultReadContainer.User,
		Password:      DefaultReadContainer.Password,
		SSLMode:       "disable",
		MaxOpenConns:  10,
		MaxIdleConns:  5,
		MaxLifetime:   10,
		Migrate:       true,
		MigrationPath: MigrationPath,
	}

	DefaultReadDb, err = connector.NewPostgres(ctx, DefaultConfig)
	if err != nil {
		panic(err)
	}

	if err := testcontainer.SeedPostgres(DefaultReadDb, DefaultSeedPath); err != nil {
		panic(err)
	}

	m.Run()
}
