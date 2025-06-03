package postgres

import (
	"astigo/internal/domain/repository"
	"astigo/pkg/dto"
	"context"
	"database/sql"
	"encoding/json"
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
        SELECT foo_id, foo.label,
    	COALESCE(
			json_agg(bar_id) FILTER (WHERE bar_id IS NOT NULL),
			'[]'
		) AS bar_ids
        FROM foo
        LEFT JOIN bar USING(foo_id)
        GROUP BY foo_id, foo.label
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
		var barIDsJSON []byte

		if err := rows.Scan(&foo.Id, &foo.Label, &barIDsJSON); err != nil {
			return nil, fmt.Errorf("error scanning foo row: %w", err)
		}

		err = json.Unmarshal(barIDsJSON, &foo.Bars)
		if err != nil {
			return nil, fmt.Errorf("error unmarshal bar id: %w", err)
		}

		foos = append(foos, foo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foo rows: %w", err)
	}

	return foos, nil
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
