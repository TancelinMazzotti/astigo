package redis

import (
	"astigo/internal/domain/cache"
	"github.com/redis/go-redis/v9"
)

var (
	_ cache.IFooCahe = (*FooRedis)(nil)
)

type FooRedis struct {
	db *redis.Client
}

func NewFooRedis(db *redis.Client) *FooRedis {
	return &FooRedis{
		db: db,
	}
}
