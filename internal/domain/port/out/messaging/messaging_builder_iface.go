package messaging

import "context"

type IMessagingBuilder interface {
	Build(ctx context.Context) error
	GetFooMessaging() (IFooMessaging, error)
}
