package postgres

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	_ repository.IFooRepository = (*FooPostgres)(nil)
)

type FooPostgres struct {
	db *sql.DB
}

func (f FooPostgres) FindAll(ctx context.Context, pagination handler.PaginationInput) ([]model.Foo, error) {
	query := `
        SELECT foo_id, foo.label, foo.secret
        FROM foo
        ORDER BY foo_id
        LIMIT $1 OFFSET $2`

	rows, err := f.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("error querying foos: %w", err)
	}
	defer rows.Close()

	var foos []model.Foo
	for rows.Next() {
		var foo model.Foo

		if err := rows.Scan(&foo.Id, &foo.Label, &foo.Secret); err != nil {
			return nil, fmt.Errorf("error scanning foo row: %w", err)
		}

		foos = append(foos, foo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foo rows: %w", err)
	}

	return foos, nil
}

func (f FooPostgres) FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	query := `
        SELECT foo_id, foo.label, foo.secret
        FROM foo
        WHERE foo_id = $1`

	row := f.db.QueryRowContext(ctx, query, id)
	var foo model.Foo

	if err := row.Scan(&foo.Id, &foo.Label, &foo.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewNotFound("foo", fmt.Sprintf("id: %s", id))
		}
		return nil, fmt.Errorf("error scanning foo row: %w", err)
	}

	return &foo, nil
}

func (f FooPostgres) Create(ctx context.Context, foo model.Foo) error {
	query := `INSERT INTO foo (foo_id,label, secret) VALUES ($1, $2, $3)`

	result, err := f.db.ExecContext(ctx, query, foo.Id, foo.Label, foo.Secret)
	if err != nil {
		return fmt.Errorf("error inserting foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (f FooPostgres) Update(ctx context.Context, foo model.Foo) error {
	query := `UPDATE foo SET label = $1, secret = $2 WHERE foo_id = $3`

	result, err := f.db.ExecContext(ctx, query, foo.Label, foo.Secret, foo.Id)
	if err != nil {
		return fmt.Errorf("error updating foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (f FooPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if err := f.DeleteBars(ctx, id); err != nil {
		return fmt.Errorf("error deleting bars: %w", err)
	}

	query := `DELETE FROM foo WHERE foo_id = $1`

	result, err := f.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (f FooPostgres) DeleteBars(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM bar WHERE foo_id = $1`

	_, err := f.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting foo: %w", err)
	}

	return nil
}

func NewFooPostgres(db *sql.DB) *FooPostgres {
	return &FooPostgres{db: db}
}
