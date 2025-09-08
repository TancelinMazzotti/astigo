package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/data"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/repository"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/repository/postgres/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
)

var (
	_ repository.IFooRepository = (*FooPostgres)(nil)
)

// FooPostgres is a concrete implementation of the IFooRepository interface that interacts with a PostgreSQL database.
type FooPostgres struct {
	db *sql.DB
}

// FindAll retrieves a list of Foo records from the database based on the provided pagination input (limit and offset).
func (f FooPostgres) FindAll(ctx context.Context, input data.FooReadListInput) ([]*model.Foo, error) {
	tracer := otel.Tracer("FooPostgres")
	ctx, span := tracer.Start(ctx, "FooPostgres.FindAll")
	defer span.End()

	span.SetAttributes(
		attribute.Int("offset", input.Offset),
		attribute.Int("limit", input.Limit),
	)

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
		span.RecordError(err)
		span.SetStatus(codes.Error, "error querying foos")
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
			span.RecordError(err)
			span.SetStatus(codes.Error, "error scanning foo row")
			return nil, fmt.Errorf("error scanning foo row: %w", err)
		}

		foo := fooEntity.ToModel()
		foos = append(foos, foo)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error iterating foo rows")
		return nil, fmt.Errorf("error iterating foo rows: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Int("result.count", len(foos)))
	return foos, nil
}

// FindByID retrieves a Foo record by its unique identifier from the database.
func (f FooPostgres) FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	tracer := otel.Tracer("FooPostgres")
	ctx, span := tracer.Start(ctx, "FooPostgres.FindByID")
	defer span.End()

	span.SetAttributes(attribute.String("foo.id", id.String()))

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
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "foo not found")
			return nil, port.NewErrNotFound("foo", "id", id.String())
		}
		span.SetStatus(codes.Error, "error scanning foo row")
		return nil, fmt.Errorf("error scanning foo row: %w", err)
	}

	foo := fooEntity.ToModel()
	span.SetStatus(codes.Ok, "")
	span.SetAttributes(
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)
	return foo, nil
}

// Create inserts a new Foo record into the database
func (f FooPostgres) Create(ctx context.Context, foo *model.Foo) error {
	tracer := otel.Tracer("FooPostgres")
	ctx, span := tracer.Start(ctx, "FooPostgres.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo.id", foo.Id.String()),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)

	query := `
    INSERT INTO foo (foo_id,label, secret, value, weight)
    VALUES ($1, $2, $3, $4, $5)
    `

	result, err := f.db.ExecContext(ctx, query, foo.Id, foo.Label, foo.Secret, foo.Value, foo.Weight)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error inserting foo")
		return fmt.Errorf("error inserting foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error getting affected rows")
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		span.SetStatus(codes.Error, "no row affected")
		return fmt.Errorf("no row affected")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (f FooPostgres) Update(ctx context.Context, foo *model.Foo) error {
	tracer := otel.Tracer("FooPostgres")
	ctx, span := tracer.Start(ctx, "FooPostgres.Update")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo.id", foo.Id.String()),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)

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
		span.RecordError(err)
		span.SetStatus(codes.Error, "error updating foo")
		return fmt.Errorf("error updating foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error getting affected rows")
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		span.SetStatus(codes.Error, "no row affected")
		return fmt.Errorf("no row affected")
	}

	foo.UpdatedAt = &now
	span.SetStatus(codes.Ok, "")
	return nil
}

func (f FooPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("FooPostgres")
	ctx, span := tracer.Start(ctx, "FooPostgres.DeleteByID")
	defer span.End()

	span.SetAttributes(attribute.String("foo.id", id.String()))

	if err := f.DeleteBars(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error deleting bars")
		return fmt.Errorf("error deleting bars: %w", err)
	}

	query := `DELETE FROM foo WHERE foo_id = $1`

	result, err := f.db.ExecContext(ctx, query, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error deleting foo")
		return fmt.Errorf("error deleting foo: %w", err)
	}

	if affectedRow, err := result.RowsAffected(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error getting affected rows")
		return fmt.Errorf("error getting affected rows: %w", err)
	} else if affectedRow == 0 {
		span.SetStatus(codes.Error, "no row affected")
		return fmt.Errorf("no row affected")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (f FooPostgres) DeleteBars(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("FooPostgres")
	ctx, span := tracer.Start(ctx, "FooPostgres.DeleteBars")
	defer span.End()

	span.SetAttributes(attribute.String("foo.id", id.String()))

	query := `DELETE FROM bar WHERE foo_id = $1`

	_, err := f.db.ExecContext(ctx, query, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error deleting bars")
		return fmt.Errorf("error deleting bars: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func NewFooPostgres(db *sql.DB) *FooPostgres {
	return &FooPostgres{db: db}
}
