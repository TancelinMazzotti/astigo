package redis

import (
	"astigo/internal/domain/repository"
	"astigo/pkg/dto"
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

func (f *FooRedis) FindAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error) {
	return f.repo.FindAll(ctx, pagination)
}

func (f *FooRedis) FindByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error) {
	return f.repo.FindByID(ctx, id)
}

func (f *FooRedis) Create(ctx context.Context, foo dto.FooRequestCreateDto) error {
	return f.repo.Create(ctx, foo)
}

func (f *FooRedis) Update(ctx context.Context, foo dto.FooRequestUpdateDto) error {
	return f.repo.Update(ctx, foo)
}

func (f *FooRedis) DeleteByID(ctx context.Context, id int) error {
	return f.repo.DeleteByID(ctx, id)
}

func NewFooRedis(repo repository.IFooRepository, db *redis.Client) *FooRedis {
	return &FooRedis{
		repo: repo,
		db:   db,
	}
}
