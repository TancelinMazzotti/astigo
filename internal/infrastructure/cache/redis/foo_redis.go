package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/cache/redis/entity"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type FooRedis struct {
	db *redis.Client
}

func (f FooRedis) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	tracer := otel.Tracer("FooRedis")
	ctx, span := tracer.Start(ctx, "FooRedis.GetByID")
	defer span.End()

	key := entity.FooKey{Id: id}
	span.SetAttributes(
		attribute.String("foo.id", id.String()),
		attribute.String("redis.key", key.GetKey()),
	)

	value, err := f.db.Get(ctx, key.GetKey()).Result()
	if errors.Is(err, redis.Nil) {
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(attribute.Bool("cache.miss", true))
		return nil, nil
	} else if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get from redis")
		return nil, fmt.Errorf("fail to find foo by id: %w", err)
	}

	var fooEntity entity.FooEntity
	if err := json.Unmarshal([]byte(value), &fooEntity); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unmarshal foo")
		return nil, fmt.Errorf("fail to unmarshal foo: %w", err)
	}

	foo := fooEntity.ToModel()
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(
		attribute.Bool("cache.hit", true),
		attribute.Int("value.size", len(value)),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)
	return foo, nil
}

func (f FooRedis) Set(ctx context.Context, foo *model.Foo, expiration time.Duration) error {
	tracer := otel.Tracer("FooRedis")
	ctx, span := tracer.Start(ctx, "FooRedis.Set")
	defer span.End()

	key := entity.FooKey{Id: foo.Id}
	span.SetAttributes(
		attribute.String("foo.id", foo.Id.String()),
		attribute.String("redis.key", key.GetKey()),
		attribute.Int64("redis.expiration", int64(expiration.Seconds())),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)

	value := entity.NewFooEntity(foo)
	valueByte, err := json.Marshal(value)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal foo")
		return fmt.Errorf("fail to marshal foo: %w", err)
	}

	span.SetAttributes(attribute.Int("value.size", len(valueByte)))

	if result := f.db.Set(ctx, key.GetKey(), valueByte, expiration); result.Err() != nil {
		span.RecordError(result.Err())
		span.SetStatus(codes.Error, "failed to set in redis")
		return fmt.Errorf("fail to set foo: %w", result.Err())
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (f FooRedis) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("FooRedis")
	ctx, span := tracer.Start(ctx, "FooRedis.DeleteByID")
	defer span.End()

	key := entity.FooKey{Id: id}
	span.SetAttributes(
		attribute.String("foo.id", id.String()),
		attribute.String("redis.key", key.GetKey()),
	)

	result := f.db.Del(ctx, key.GetKey())
	if result.Err() != nil {
		span.RecordError(result.Err())
		span.SetStatus(codes.Error, "failed to delete from redis")
		return fmt.Errorf("fail to delete foo: %w", result.Err())
	}

	span.SetAttributes(attribute.Int64("redis.deleted_count", result.Val()))
	span.SetStatus(codes.Ok, "")
	return nil
}

func NewFooRedis(db *redis.Client) *FooRedis {
	return &FooRedis{db: db}
}
