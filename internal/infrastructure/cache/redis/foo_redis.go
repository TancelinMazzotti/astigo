package redis

import (
	"astigo/internal/domain/cache"
	"astigo/internal/domain/model"
	"astigo/internal/infrastructure/cache/redis/entity"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	_ cache.IFooCache = (*FooRedis)(nil)
)

type FooRedis struct {
	db *redis.Client
}

func (f FooRedis) GetByID(ctx context.Context, id int) (*model.Foo, error) {
	key := entity.FooKey{Id: id}

	value, err := f.db.Get(ctx, key.GetKey()).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("fail to find foo by id: %w", err)
	}

	var foo model.Foo
	if err := json.Unmarshal([]byte(value), &foo); err != nil {
		return nil, fmt.Errorf("fail to unmarshal foo: %w", err)
	}

	return &foo, nil
}

func (f FooRedis) Set(ctx context.Context, foo model.Foo, expiration time.Duration) error {
	key := entity.FooKey{Id: foo.Id}

	value, err := json.Marshal(foo)
	if err != nil {
		return fmt.Errorf("fail to marshal foo: %w", err)
	}

	if result := f.db.Set(ctx, key.GetKey(), value, expiration); result.Err() != nil {
		return fmt.Errorf("fail to set foo: %w", result.Err())
	}

	return nil
}

func NewFooRedis(db *redis.Client) *FooRedis {
	return &FooRedis{
		db: db,
	}
}
