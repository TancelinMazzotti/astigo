package redis

import (
	"context"
	"strings"
	"testing"

	"github.com/TancelinMazzotti/astigo/internal/infrastructure/connector"
	"github.com/TancelinMazzotti/astigo/test/helpers"
	"github.com/TancelinMazzotti/astigo/test/testcontainer"
	"github.com/redis/go-redis/v9"
)

var (
	ConfigPath      string
	DefaultSeedPath string

	DefaultConfig    connector.RedisConfig
	DefaultContainer *testcontainer.RedisContainer

	DefaultClient *redis.Client
)

func TestMain(m *testing.M) {
	var err error
	ConfigPath = strings.Join([]string{helpers.ProjectRoot, "/test/testdata/redis/init-db.conf"}, "")
	DefaultSeedPath = strings.Join([]string{helpers.ProjectRoot, "/test/testdata/redis/seed.json"}, "")

	ctx := context.Background()
	DefaultContainer, err = testcontainer.CreateRedisContainer(ctx, ConfigPath)
	if err != nil {
		panic(err)
	}

	DefaultConfig = connector.RedisConfig{
		Host:     DefaultContainer.Config.Host,
		Port:     DefaultContainer.Config.Port,
		Password: "redis",
		DB:       0,
	}

	DefaultClient, err = connector.NewRedis(ctx, DefaultConfig)
	if err != nil {
		panic(err)
	}

	if err := testcontainer.SeedRedis(ctx, DefaultClient, DefaultSeedPath); err != nil {
		panic(err)
	}

	m.Run()
}
