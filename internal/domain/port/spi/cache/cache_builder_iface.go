package cache

import "context"

type ICacheBuilder interface {
	Build(ctx context.Context) error
	GetFooCache() (IFooCache, error)
}
