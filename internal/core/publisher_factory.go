package core

import (
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/adapter/out/publisher/nats"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/publisher"
)

func GetPublisherBuilder(config Config) (publisher.IPublisherBuilder, error) {
	switch config.Driver.Publisher {
	case nats.Driver:
		return nats.NewPublisherBuilder(config.Nats), nil
	default:
		return nil, fmt.Errorf("unknown publisher driver: %s", config.Driver.Publisher)
	}
}
