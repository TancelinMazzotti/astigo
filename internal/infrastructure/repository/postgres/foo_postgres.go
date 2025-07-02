package postgres

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"astigo/internal/domain/repository"
	"astigo/internal/infrastructure/repository/postgres/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var (
	_ repository.IFooRepository = (*FooPostgres)(nil)
)

// FooPostgres is a concrete implementation of the IFooRepository interface that interacts with a PostgreSQL database.
type FooPostgres struct {
	db *sql.DB
}

// FindAll retrieves a list of Foo records from the database based on the provided pagination input (limit and offset).
func (f FooPostgres) FindAll(ctx context.Context, input handler.FooReadListInput) ([]*model.Foo, error) {
	query := `
        SELECT 
            foo.foo_id,
            foo.label,
            foo.secret,
            foo.value,
            foo.weight,
            foo.created_at,
            foo.updated_at
        FROM foo
        ORDER BY foo.foo_id
        LIMIT $1 OFFSET $2`

	rows, err := f.db.QueryContext(ctx, query, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("error querying foos: %w", err)
	}
	defer rows.Close()

	var foos []*model.Foo
	for rows.Next() {
		fooEntity := entity.Foo{}
		if err := rows.Scan(
			&fooEntity.FooId,
			&fooEntity.Label,
			&fooEntity.Secret,
			&fooEntity.Value,
			&fooEntity.Weight,
			&fooEntity.CreatedAt,
			&fooEntity.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning foo row: %w", err)
		}

		foo := fooEntity.ToModel()
		foos = append(foos, foo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating foo rows: %w", err)
	}

	return foos, nil
}

// FindByID retrieves a Foo record by its unique identifier from the database. Returns the Foo model or an error if not found.
func (f FooPostgres) FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	query := `
        SELECT
            foo.foo_id,
            foo.label,
            foo.secret,
            foo.value,
            foo.weight,
            foo.created_at,
            foo.updated_at
        FROM foo
        WHERE foo.foo_id = $1`

	row := f.db.QueryRowContext(ctx, query, id)

	fooEntity := entity.Foo{}
	if err := row.Scan(
		&fooEntity.FooId,
		&fooEntity.Label,
		&fooEntity.Secret,
		&fooEntity.Value,
		&fooEntity.Weight,
		&fooEntity.CreatedAt,
		&fooEntity.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewNotFound("foo", fmt.Sprintf("id: %s", id))
		}
		return nil, fmt.Errorf("error scanning foo row: %w", err)
	}

	foo := fooEntity.ToModel()

	return foo, nil
}

// Create inserts a new Foo record into the database and returns an error if the operation fails.
func (f FooPostgres) Create(ctx context.Context, foo *model.Foo) error {
	query := `
	INSERT INTO foo (foo_id,label, secret, value, weight)
	VALUES ($1, $2, $3, $4, $5)
	`

	result, err := f.db.ExecContext(ctx, query, foo.Id, foo.Label, foo.Secret, foo.Value, foo.Weight)
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

func (f FooPostgres) Update(ctx context.Context, foo *model.Foo) error {
	now := time.Now()
	query := `
	UPDATE foo 
	SET label = $1, 
	    secret = $2,
	    value = $3,
	    weight = $4,
	    updated_at = $5
	WHERE foo_id = $6
	`

	result, err := f.db.ExecContext(ctx, query, foo.Label, foo.Secret, foo.Value, foo.Weight, now, foo.Id)
	if err != nil {
		return fmt.Errorf("error updating foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		return fmt.Errorf("no row affected")
	}

	foo.UpdatedAt = &now

	return nil
}

// DeleteByID removes a Foo record and its associated Bar records from the database by the provided unique identifier.
// Returns an error if the operation fails or no records are affected.
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

// DeleteBars removes all Bar records associated with a given Foo ID from the database. Returns an error if the operation fails.
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
