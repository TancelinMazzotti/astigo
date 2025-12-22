package testcontainer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/TancelinMazzotti/astigo/internal/adapter/out/connector"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// RedisContainer represents a wrapper around a Redis container with configuration and container-related functionalities.
type RedisContainer struct {
	*redis.RedisContainer
	Config connector.RedisConfig
}

// CreateRedisContainer initializes and runs a Redis container with specified configurations and seeds data from JSON.
// Returns a RedisContainer instance and any error encountered during creation or data seeding.
func CreateRedisContainer(ctx context.Context, pathConfig string) (*RedisContainer, error) {
	config := connector.RedisConfig{
		Host:     "localhost",
		Password: "redis",
		DB:       0,
	}

	redisContainer, err := redis.Run(ctx, "redis:6.2.6-alpine",
		redis.WithConfigFile(pathConfig),
	)
	if err != nil {
		return nil, err
	}

	containerPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, err
	}
	config.Port = containerPort.Int()

	return &RedisContainer{
		RedisContainer: redisContainer,
		Config:         config,
	}, nil
}

// SeedRedis reads a JSON file and populates a Redis instance with key-value pairs derived from the file's content.
func SeedRedis(ctx context.Context, client *redis2.Client, seedPath string) error {
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var entries map[string]interface{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for key, value := range entries {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		err = client.Set(ctx, key, jsonValue, 0).Err()
		if err != nil {
			return fmt.Errorf("failed to set key %s: %w", key, err)
		}
	}

	return nil
}
