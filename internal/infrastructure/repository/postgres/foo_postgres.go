package postgres

import (
	"astigo/internal/domain/repository"
	"astigo/pkg/dto"
	"context"
	"database/sql"
	"fmt"
)

var (
	_ repository.IFooRepository = (*FooPostgres)(nil)
)

type FooPostgres struct {
	db *sql.DB
}

func (f FooPostgres) FindAll(ctx context.Context, pagination dto.PaginationRequestDto) ([]dto.FooResponseReadDto, error) {
	query := `
        SELECT foo_id, label
        FROM foo
        ORDER BY foo_id
        LIMIT $1 OFFSET $2`

	rows, err := f.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("error querying foos: %w", err)
	}
	defer rows.Close()

	var foos []dto.FooResponseReadDto
	for rows.Next() {
		var foo dto.FooResponseReadDto

		if err := rows.Scan(&foo.Id, &foo.Label); err != nil {
			return nil, fmt.Errorf("error scanning foo row: %w", err)
		}
		foos = append(foos, foo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foo rows: %w", err)
	}

	for i, foo := range foos {
		bars, err := f.findAllBarId(ctx, foo.Id)
		if err != nil {
			return nil, fmt.Errorf("error querying bars: %w", err)
		}
		foos[i].Bars = bars
	}

	return foos, nil
}

func (f FooPostgres) findAllBarId(ctx context.Context, fooID int) ([]int, error) {
	query := "SELECT bar_id FROM bar JOIN foo USING(foo_id) WHERE foo_id = $1"
	rows, err := f.db.QueryContext(ctx, query, fooID)
	if err != nil {
		return nil, fmt.Errorf("error querying foos: %w", err)
	}
	defer rows.Close()

	var barsID []int
	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning foo row: %w", err)
		}
		barsID = append(barsID, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foo rows: %w", err)
	}

	return barsID, nil
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
