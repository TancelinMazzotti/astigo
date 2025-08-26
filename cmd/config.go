package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set
func initConfig() error {
	setDefaults()

	// Config file setup
	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("fatal error reading configuration file: %w", err)
		}
		log.Println("No configuration file found, using defaults and environment variables")
	} else {
		log.Printf("Using configuration file: %s", viper.ConfigFileUsed())
	}

	// Enable environment variable support
	viper.SetEnvPrefix("ASTIGO")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return nil
}

// setDefaults initializes default configuration values for all components
func setDefaults() {
	// HTTP server defaults
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.mode", "debug")

	// gRPC server defaults
	viper.SetDefault("grpc.port", 50051)

	// Auth configuration defaults
	viper.SetDefault("auth.issuer", "http://localhost:8080/realms/astigo")
	viper.SetDefault("auth.client_id", "astigo-api")

	// Logging configuration defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.encoding", "json")

	// Jaeger configuration
	viper.SetDefault("jaeger.url", "localhost:4318")
	viper.SetDefault("jaeger.service_name", "astigo")

	// PostgreSQL connection defaults
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("postgres.max_open_conns", 10)
	viper.SetDefault("postgres.max_idle_conns", 5)
	viper.SetDefault("postgres.max_lifetime", 300)

	// Redis connection defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// NATS connection defaults
	viper.SetDefault("nats.url", "nats://localhost:4222")
}
