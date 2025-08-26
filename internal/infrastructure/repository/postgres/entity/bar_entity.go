package entity

import (
	"database/sql"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"

	"github.com/google/uuid"
)

// Bar represents a database model with nullable fields for bar entities.
type Bar struct {
	BarId     sql.Null[uuid.UUID] `db:"bar_id"`
	Label     sql.NullString      `db:"label"`
	Secret    sql.NullString      `db:"secret"`
	Value     sql.NullInt32       `db:"value"`
	FooId     sql.Null[uuid.UUID] `db:"foo_id"`
	CreatedAt sql.NullTime        `db:"created_at"`
	UpdatedAt sql.NullTime        `db:"updated_at"`
}

// ToModel converts a database model of Bar into a domain-level model.Bar instance with non-nullable fields.
func (b *Bar) ToModel() *model.Bar {
	bar := model.Bar{}
	if b.BarId.Valid {
		bar.Id = b.BarId.V
	}
	if b.Label.Valid {
		bar.Label = b.Label.String
	}
	if b.Secret.Valid {
		bar.Secret = b.Secret.String
	}
	if b.Value.Valid {
		bar.Value = int(b.Value.Int32)
	}
	if b.FooId.Valid {
		bar.FooID = b.FooId.V
	}
	if b.CreatedAt.Valid {
		bar.CreatedAt = b.CreatedAt.Time
	}
	if b.UpdatedAt.Valid {
		bar.UpdatedAt = &b.UpdatedAt.Time
	}

	return &bar
}
