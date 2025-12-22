package core

import (
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/messaging"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/messaging/nats"
)

func GetMessagingBuilder(config Config) (messaging.IMessagingBuilder, error) {
	switch config.Driver.Publisher {
	case nats.Driver:
		return nats.NewMessagingBuilder(config.Nats), nil
	default:
		return nil, fmt.Errorf("unknown messaging driver: %s", config.Driver.Publisher)
	}
}
