package dto

import "github.com/google/uuid"

type FooReadRequest struct {
	Id string `uri:"id" binding:"required"`
}

type FooReadResponse struct {
	Id    uuid.UUID `json:"id" binding:"required"`
	Label string    `json:"label" binding:"required"`
}

type FooCreateBody struct {
	Label  string `json:"label" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type FooUpdateRequest struct {
	Id string `uri:"id" binding:"required"`
}
type FooUpdateBody struct {
	Label  string `json:"label" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type FooDeleteRequest struct {
	Id string `uri:"id" binding:"required"`
}
