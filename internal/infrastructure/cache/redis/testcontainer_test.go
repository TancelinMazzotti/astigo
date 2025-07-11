package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"os"
)

// RedisContainer represents a wrapper around a Redis container with configuration and container-related functionalities.
type RedisContainer struct {
	*redis.RedisContainer
	Config RedisConfig
}

// CreateRedisContainer initializes and runs a Redis container with specified configurations and seeds data from JSON.
// Returns a RedisContainer instance and any error encountered during creation or data seeding.
func CreateRedisContainer(ctx context.Context) (*RedisContainer, error) {
	config := RedisConfig{
		Host:     "localhost",
		Password: "redis",
		DB:       0,
	}

	redisContainer, err := redis.Run(ctx, "redis:6.2.6-alpine",
		redis.WithConfigFile("init-db.conf"),
	)
	if err != nil {
		return nil, err
	}

	containerPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, err
	}
	config.Port = containerPort.Int()

	if err := SeedFromJSON(ctx, config); err != nil {
		return nil, err
	}

	return &RedisContainer{
		RedisContainer: redisContainer,
		Config:         config,
	}, nil
}

// SeedFromJSON reads a JSON file and populates a Redis instance with key-value pairs derived from the file's content.
func SeedFromJSON(ctx context.Context, config RedisConfig) error {
	data, err := os.ReadFile("testdata.json")
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var entries map[string]interface{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	r, err := NewRedis(ctx, config)
	if err != nil {
		return err
	}

	for key, value := range entries {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		err = r.Set(ctx, key, jsonValue, 0).Err()
		if err != nil {
			return fmt.Errorf("failed to set key %s: %w", key, err)
		}
	}

	return nil
}
