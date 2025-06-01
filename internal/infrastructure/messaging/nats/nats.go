package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type NatsConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func NewNats(config NatsConfig) (*nats.Conn, error) {
	conn, err := nats.Connect(config.URL, nats.UserInfo(config.Username, config.Password))
	if err != nil {
		return nil, fmt.Errorf("can't open connection to Nats: %w", err)
	}

	return conn, nil
}
