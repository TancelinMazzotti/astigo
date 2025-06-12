package config

import (
	"astigo/internal/application/http"
	"astigo/internal/infrastructure/cache/redis"
	"astigo/internal/infrastructure/messaging/nats"
	"astigo/internal/infrastructure/repository/postgres"
	"astigo/internal/tool"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	Gin http.GinConfig    `mapstructure:"http"`
	Log tool.LoggerConfig `mapstructure:"log"`

	Postgres postgres.PostgresConfig `mapstructure:"postgres"`
	Nats     nats.NatsConfig         `mapstructure:"nats"`
	Redis    redis.RedisConfig       `mapstructure:"redis"`
}

func Load() error {
	// 1. Configuration par défaut
	setDefaults()

	// 2. Chargement du fichier de configuration
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("info: impossible de lire le fichier de config: %v\n", err)
	}

	// 3. Configuration des variables d'environnement
	viper.SetEnvPrefix("ASTIGO")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("erreur parsing config: %v\n", err)
	}

	Cfg = &cfg
	return nil
}

func setDefaults() {
	// Postgres defaults
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("postgres.max_open_conns", 10)
	viper.SetDefault("postgres.max_idle_conns", 5)
	viper.SetDefault("postgres.max_lifetime", 300)

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// HTTP defaults
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.mode", "debug")

	// NATS defaults
	viper.SetDefault("nats.url", "nats://localhost:4222")

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.encoding", "json")
}
