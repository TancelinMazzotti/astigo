package publisher

import "context"

type IPublisherBuilder interface {
	Build(ctx context.Context) error
	GetFooPublisher() (IFooPublisher, error)
}
