package core

import (
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/publisher"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/publisher/nats"
)

func GetPublisherBuilder(config Config) (publisher.IPublisherBuilder, error) {
	switch config.Driver.Publisher {
	case nats.Driver:
		return nats.NewPublisherBuilder(config.Nats), nil
	default:
		return nil, fmt.Errorf("unknown publisher driver: %s", config.Driver.Publisher)
	}
}
