package core

import (
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/adapter/out/repository/postgres"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/repository"
)

func GetRepositoryBuilder(config Config) (repository.IRepositoryBuilder, error) {
	switch config.Driver.Repository {
	case postgres.Driver:
		return postgres.NewRepositoryBuilder(config.Postgres), nil
	default:
		return nil, fmt.Errorf("unknown repository driver: %s", config.Driver.Repository)
	}
}
