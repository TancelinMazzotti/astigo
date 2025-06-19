package model

import "github.com/google/uuid"

type Bar struct {
	Id     uuid.UUID `json:"id"`
	Label  string    `json:"label"`
	Secret string    `json:"secret"`
}
