package core

import (
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/cache"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/cache/redis"
)

func GetCacheBuilder(config Config) (cache.ICacheBuilder, error) {
	switch config.Driver.Cache {
	case redis.Driver:
		return redis.NewCacheBuilder(config.Redis), nil
	default:
		return nil, fmt.Errorf("unknown cache driver: %s", config.Driver.Cache)
	}
}
