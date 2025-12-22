package repository

import "context"

type IRepositoryBuilder interface {
	Build(ctx context.Context) error
	GetFooRepository() (IFooRepository, error)
}
