package redis

import (
	"astigo/internal/domain/repository"
	"astigo/pkg/model"
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	_ repository.IFooRepository = (*FooRedis)(nil)
)

type FooRedis struct {
	repo repository.IFooRepository
	db   *redis.Client
}

func NewFooRedis(repo repository.IFooRepository, db *redis.Client) *FooRedis {
	return &FooRedis{
		repo: repo,
		db:   db,
	}
}

func (f *FooRedis) FindByID(ctx context.Context, id string) (*model.Foo, error) {
	return f.repo.FindByID(ctx, id)
}

func (f *FooRedis) Create(ctx context.Context, foo *model.Foo) error {
	return f.repo.Create(ctx, foo)
}
