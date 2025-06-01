package postgres

import (
	"astigo/internal/domain/repository"
	"astigo/internal/tool"
	"astigo/pkg/model"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

var (
	_ repository.IFooRepository = (*FooPostgres)(nil)
)

type FooPostgres struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *FooPostgres {
	return &FooPostgres{db: db}
}

func (f *FooPostgres) FindByID(ctx context.Context, id string) (*model.Foo, error) {
	const query = `SELECT id, name FROM foo WHERE id = $1`
	row := f.db.QueryRowContext(ctx, query, id)

	var foo model.Foo
	if err := row.Scan(&foo.Id, &foo.Label); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tool.Logger.Debug("Foo not found", zap.String("id", id))
			return nil, err
		}
		tool.Logger.Error("Error scanning foo", zap.Error(err))
		return nil, err
	}

	return &foo, nil
}

func (f *FooPostgres) Create(ctx context.Context, foo *model.Foo) error {
	const query = `INSERT INTO foo (id, name) VALUES ($1, $2)`
	_, err := f.db.ExecContext(ctx, query, foo.Id, foo.Label)
	if err != nil {
		tool.Logger.Error("Error creating foo",
			zap.Int("id", foo.Id),
			zap.Error(err))
	}
	return err
}
