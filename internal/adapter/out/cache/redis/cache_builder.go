package redis

import (
	"context"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/adapter/out/connector"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/cache"
	"github.com/redis/go-redis/v9"
)

const Driver = "redis"

var (
	_ cache.ICacheBuilder = (*CacheBuilder)(nil)
)

type CacheBuilder struct {
	config connector.RedisConfig
	client *redis.Client
	foo    *FooRedis
}

func (c *CacheBuilder) Build(ctx context.Context) error {
	var err error

	c.client, err = connector.NewRedis(ctx, c.config)
	if err != nil {
		return fmt.Errorf("failed to create redis client: %w", err)
	}

	c.foo = NewFooRedis(c.client)

	return nil
}

func (c *CacheBuilder) GetFooCache() (cache.IFooCache, error) {
	if c.foo == nil {
		return nil, fmt.Errorf("foo cache is not initialized")
	}
	return c.foo, nil
}

func NewCacheBuilder(config connector.RedisConfig) cache.ICacheBuilder {
	return &CacheBuilder{config: config}
}
