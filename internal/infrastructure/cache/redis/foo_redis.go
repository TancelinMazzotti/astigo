package redis

import (
	"astigo/internal/domain/contract/cache"
	"astigo/internal/domain/model"
	"astigo/internal/infrastructure/cache/redis/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	_ cache.IFooCache = (*FooRedis)(nil)
)

// FooRedis is a type that provides Redis-based operations for managing Foo entities.
type FooRedis struct {
	db *redis.Client
}

// GetByID retrieves a Foo entity by its UUID from Redis and converts it to a model.Foo. Returns nil if not found.
func (f FooRedis) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	key := entity.FooKey{Id: id}

	value, err := f.db.Get(ctx, key.GetKey()).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("fail to find foo by id: %w", err)
	}

	var fooEntity entity.FooEntity
	if err := json.Unmarshal([]byte(value), &fooEntity); err != nil {
		return nil, fmt.Errorf("fail to unmarshal foo: %w", err)
	}

	foo := fooEntity.ToModel()

	return foo, nil

}

// Set stores a given Foo model in Redis with an optional expiration duration. Returns an error if the operation fails.
func (f FooRedis) Set(ctx context.Context, foo *model.Foo, expiration time.Duration) error {
	key := entity.FooKey{Id: foo.Id}
	value := entity.NewFooEntity(foo)

	valueByte, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("fail to marshal foo: %w", err)
	}

	if result := f.db.Set(ctx, key.GetKey(), valueByte, expiration); result.Err() != nil {
		return fmt.Errorf("fail to set foo: %w", result.Err())
	}

	return nil
}

// DeleteByID removes a Foo entity from Redis using its UUID. Returns an error if the deletion operation fails.
func (f FooRedis) DeleteByID(ctx context.Context, id uuid.UUID) error {
	key := entity.FooKey{Id: id}

	if result := f.db.Del(ctx, key.GetKey()); result.Err() != nil {
		return fmt.Errorf("fail to delete foo: %w", result.Err())
	}

	return nil
}

func NewFooRedis(db *redis.Client) *FooRedis {
	return &FooRedis{
		db: db,
	}
}
