package redis

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

var (
	_ cache.IFooCahe = (*FooRedis)(nil)
)

type FooRedis struct {
	repo repository.IFooRepository
	db   *redis.Client
}

func NewFooRedis(db *redis.Client) *FooRedis {
	return &FooRedis{
		db: db,
	}
}
