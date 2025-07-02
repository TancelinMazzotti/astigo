package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// RedisConfig represents the configuration settings required to connect to a Redis server.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NewRedis initializes a new Redis client using the provided configuration and verifies the connection with a ping.
func NewRedis(config RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Vérifier la connexion
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("échec de connexion à Redis : %w", err)
	}

	return client, nil
}
