package core

import (
	"github.com/TancelinMazzotti/astigo/internal/adapter/in/grpc"
	"github.com/TancelinMazzotti/astigo/internal/adapter/in/http"
	"github.com/TancelinMazzotti/astigo/internal/adapter/out/storage/s3storage"
	"github.com/TancelinMazzotti/astigo/internal/adapter/out/telemetry"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/connector"
)

// Config represents the main configuration structure for the application, encompassing all required service settings.
type Config struct {
	Log LoggerConfig `mapstructure:"log"`

	Gin       http.Config      `mapstructure:"http"`
	Grpc      grpc.Config      `mapstructure:"grpc"`
	Telemetry telemetry.Config `mapstructure:"telemetry"`
	Auth      struct {
		ClientID string `mapstructure:"client_id"`
		Issuer   string `mapstructure:"issuer"`
	} `mapstructure:"auth"`

	Driver DriverConfig `mapstructure:"driver"`

	Postgres connector.PostgresConfig `mapstructure:"postgres"`
	Nats     connector.NatsConfig     `mapstructure:"nats"`
	Redis    connector.RedisConfig    `mapstructure:"redis"`
	S3       s3storage.Config         `mapstructure:"s3"`
}

type DriverConfig struct {
	Repository string `mapstructure:"repository"`
	Cache      string `mapstructure:"cache"`
	Publisher  string `mapstructure:"publisher"`
	Subscriber string `mapstructure:"subscriber"`
	Storage    string `mapstructure:"storage"`
}
