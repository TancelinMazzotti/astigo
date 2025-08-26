package entity

import (
	"database/sql"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"

	"github.com/google/uuid"
)

// Foo represents a database entity with nullable fields for handling record information like ID, label, and timestamps.
type Foo struct {
	FooId     sql.Null[uuid.UUID] `db:"foo_id"`
	Label     sql.NullString      `db:"label"`
	Secret    sql.NullString      `db:"secret"`
	Value     sql.NullInt32       `db:"value"`
	Weight    sql.NullFloat64     `db:"weight"`
	CreatedAt sql.NullTime        `db:"created_at"`
	UpdatedAt sql.NullTime        `db:"updated_at"`
}

// ToModel converts a database model of Foo into a domain-level model.Foo instance with non-nullable fields.
func (f *Foo) ToModel() *model.Foo {
	foo := model.Foo{}
	if f.FooId.Valid {
		foo.Id = f.FooId.V
	}
	if f.Label.Valid {
		foo.Label = f.Label.String
	}
	if f.Secret.Valid {
		foo.Secret = f.Secret.String
	}
	if f.Value.Valid {
		foo.Value = int(f.Value.Int32)
	}
	if f.Weight.Valid {
		foo.Weight = float32(f.Weight.Float64)
	}
	if f.CreatedAt.Valid {
		foo.CreatedAt = f.CreatedAt.Time
	}
	if f.UpdatedAt.Valid {
		foo.UpdatedAt = &f.UpdatedAt.Time
	}

	return &foo
}
