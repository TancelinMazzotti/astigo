package messaging

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
)

// IFooMessaging defines a contract for publishing events related to Foo entities.
// PublishFooCreated sends a message when a Foo entity is created.
// PublishFooUpdated sends a message when a Foo entity is updated.
// PublishFooDeleted sends a message when a Foo entity is deleted.
type IFooMessaging interface {
	PublishFooCreated(ctx context.Context, foo *model.Foo) error
	PublishFooUpdated(ctx context.Context, foo *model.Foo) error
	PublishFooDeleted(ctx context.Context, id uuid.UUID) error
}
