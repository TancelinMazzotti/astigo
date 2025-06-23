package messaging

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
)

type IFooMessaging interface {
	PublishFooCreated(ctx context.Context, foo model.Foo) error
	PublishFooUpdated(ctx context.Context, foo model.Foo) error
	PublishFooDeleted(ctx context.Context, id uuid.UUID) error
}
