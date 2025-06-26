package model

import (
	"github.com/google/uuid"
	"time"
)

type Bar struct {
	Id     uuid.UUID
	Label  string
	Secret string
	Value  int
	FooID  uuid.UUID
	Foo    *Foo

	CreatedAt time.Time
	UpdatedAt *time.Time
}
