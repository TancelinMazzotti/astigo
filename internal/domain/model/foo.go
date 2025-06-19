package model

import "github.com/google/uuid"

type Foo struct {
	Id     uuid.UUID `json:"id"`
	Label  string    `json:"label"`
	Secret string    `json:"secret"`
}
