package dto

import (
	"astigo/internal/domain/model"
	"github.com/google/uuid"
)

type FooReadRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type FooReadResponse struct {
	Id     uuid.UUID `json:"id" binding:"required"`
	Label  string    `json:"label" binding:"required"`
	Value  int       `json:"value" binding:"required"`
	Weight float32   `json:"weight" binding:"required"`
}

func NewFooReadResponse(foo *model.Foo) *FooReadResponse {
	return &FooReadResponse{
		Id:     foo.Id,
		Label:  foo.Label,
		Value:  foo.Value,
		Weight: foo.Weight,
	}
}

type FooCreateBody struct {
	Label  string  `json:"label" binding:"required"`
	Secret string  `json:"secret" binding:"required"`
	Value  int     `json:"value" binding:"required"`
	Weight float32 `json:"weight" binding:"required"`
}

type FooCreateResponse struct {
	Id uuid.UUID `json:"id" binding:"required"`
}

type FooUpdateRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}
type FooUpdateBody struct {
	Label  string  `json:"label" binding:"required"`
	Secret string  `json:"secret" binding:"required"`
	Value  int     `json:"value" binding:"required"`
	Weight float32 `json:"weight" binding:"required"`
}

type FooPatchRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}
type FooPatchBody struct {
	Label  *string  `json:"label" binding:"omitempty"`
	Secret *string  `json:"secret" binding:"omitempty"`
	Value  *int     `json:"value" binding:"omitempty"`
	Weight *float32 `json:"weight" binding:"omitempty"`
}

type FooDeleteRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}
