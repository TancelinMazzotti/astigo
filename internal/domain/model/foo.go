package model

import (
	"github.com/google/uuid"
	"time"
)

type Foo struct {
	Id     uuid.UUID
	Label  string
	Secret string
	Value  int
	Weight float32

	CreatedAt time.Time
	UpdatedAt *time.Time

	Bars []*Bar
}
