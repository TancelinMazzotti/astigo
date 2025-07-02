package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

// NatsConfig represents the configuration settings required to connect to a NATS server.
type NatsConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// NewNats establishes a new NATS connection using the provided configuration and returns the connection instance or an error.
func NewNats(config NatsConfig) (*nats.Conn, error) {
	conn, err := nats.Connect(config.URL, nats.UserInfo(config.Username, config.Password))
	if err != nil {
		return nil, fmt.Errorf("can't open connection to Nats: %w", err)
	}

	return conn, nil
}
