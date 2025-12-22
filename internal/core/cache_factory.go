package core

import (
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/adapter/out/cache/redis"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/cache"
)

func GetCacheBuilder(config Config) (cache.ICacheBuilder, error) {
	switch config.Driver.Cache {
	case redis.Driver:
		return redis.NewCacheBuilder(config.Redis), nil
	default:
		return nil, fmt.Errorf("unknown cache driver: %s", config.Driver.Cache)
	}
}
