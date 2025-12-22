package publisher

import (
	"context"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"

	"github.com/google/uuid"
)

// IFooPublisher defines a port for publishing events related to Foo entities.
// PublishCreated sends a message when a Foo entity is created.
// PublishUpdated sends a message when a Foo entity is updated.
// PublishDeleted sends a message when a Foo entity is deleted.
type IFooPublisher interface {
	PublishCreated(ctx context.Context, foo *model.Foo) error
	PublishUpdated(ctx context.Context, foo *model.Foo) error
	PublishDeleted(ctx context.Context, id uuid.UUID) error
}
