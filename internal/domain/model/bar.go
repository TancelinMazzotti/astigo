package model

import (
	"time"

	"github.com/google/uuid"
)

type Bar struct {
	Id     uuid.UUID `validate:"required"`
	Label  string    `validate:"required,min=3,max=100"`
	Secret string    `validate:"required,min=3,max=100"`
	Value  int       `validate:"required,gte=0,lte=1000"`
	FooID  uuid.UUID `validate:"required"`
	Foo    *Foo      `validate:"-"`

	CreatedAt time.Time  `validate:"omitempty,datetime"`
	UpdatedAt *time.Time `validate:"omitempty,datetime"`
}
