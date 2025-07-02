package model

import (
	"github.com/google/uuid"
	"time"
)

type Foo struct {
	Id     uuid.UUID `validate:"required"`
	Label  string    `validate:"required,min=3,max=100"`
	Secret string    `validate:"required,min=3,max=100"`
	Value  int       `validate:"required,gte=0,lte=1000"`
	Weight float32   `validate:"required,gte=0"`

	CreatedAt time.Time  `validate:"omitempty"`
	UpdatedAt *time.Time `validate:"omitempty"`

	Bars []*Bar `validate:"dive"`
}
