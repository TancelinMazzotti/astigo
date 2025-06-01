package postgres

import (
	"astigo/internal/domain/repository"
	"astigo/pkg/dto"
	"context"
	"database/sql"
)

var (
	_ repository.IFooRepository = (*FooPostgres)(nil)
)

type FooPostgres struct {
	db *sql.DB
}

func (f FooPostgres) FindAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error) {
	//TODO implement me
	panic("implement me")
}

func (f FooPostgres) FindByID(ctx context.Context, id int) (*dto.FooResponseReadDto, error) {
	//TODO implement me
	panic("implement me")
}

func (f FooPostgres) Create(ctx context.Context, foo dto.FooRequestCreateDto) error {
	//TODO implement me
	panic("implement me")
}

func (f FooPostgres) Update(ctx context.Context, foo dto.FooRequestUpdateDto) error {
	//TODO implement me
	panic("implement me")
}

func (f FooPostgres) DeleteByID(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewFooPostgres(db *sql.DB) *FooPostgres {
	return &FooPostgres{db: db}
}
