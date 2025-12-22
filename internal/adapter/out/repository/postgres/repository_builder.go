package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/repository"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/connector"
)

const Driver = "postgres"

var (
	_ repository.IRepositoryBuilder = (*RepositoryBuilder)(nil)
)

type RepositoryBuilder struct {
	config connector.PostgresConfig
	db     *sql.DB
	foo    *FooPostgres
}

func (p *RepositoryBuilder) Build(ctx context.Context) error {
	var err error

	p.db, err = connector.NewPostgres(ctx, p.config)
	if err != nil {
		return fmt.Errorf("failed to create postgres client: %w", err)
	}

	p.foo = NewFooPostgres(p.db)

	return nil
}

func (p *RepositoryBuilder) GetFooRepository() (repository.IFooRepository, error) {
	if p.foo == nil {
		return nil, fmt.Errorf("foo repository is not initialized")
	}
	return p.foo, nil
}

func NewRepositoryBuilder(config connector.PostgresConfig) repository.IRepositoryBuilder {
	return &RepositoryBuilder{config: config}
}
