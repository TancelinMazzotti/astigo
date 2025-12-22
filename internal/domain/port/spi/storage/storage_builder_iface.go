package storage

import "context"

type IStorageBuilder interface {
	Build(ctx context.Context) error
	GetFileStorage() (IFileStorage, error)
}
