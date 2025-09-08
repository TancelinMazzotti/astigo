package service

import (
	"context"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/data"

	"github.com/google/uuid"
)

// IFooService defines the interface for handling operations related to Foo entities.
// GetAll retrieves a list of Foo entities based on the provided input.
// GetByID fetches a Foo entity by its unique identifier.
// Create adds a new Foo entity based on the provided input and returns the created instance.
// Update modifies an existing Foo entity based on the provided input.
// DeleteByID removes a Foo entity identified by its unique identifier.
type IFooService interface {
	GetAll(ctx context.Context, input data.FooReadListInput) ([]*model.Foo, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Create(ctx context.Context, input data.FooCreateInput) (*model.Foo, error)
	Update(ctx context.Context, input data.IFooUpdateMerger) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
