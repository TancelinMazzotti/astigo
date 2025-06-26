package repository

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
	"github.com/google/uuid"
)

// IFooRepository represents a contract for interacting with Foo data storage.
// FindAll retrieves a paginated list of Foo entities from the repository.
// FindByID fetches a Foo entity by its unique identifier.
// Create adds a new Foo entity to the repository.
// Update modifies an existing Foo entity in the repository.
// DeleteByID removes a Foo entity by its unique identifier from the repository.
type IFooRepository interface {
	FindAll(ctx context.Context, pagination handler.FooReadListInput) ([]*model.Foo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Foo, error)
	Create(ctx context.Context, foo *model.Foo) error
	Update(ctx context.Context, foo *model.Foo) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
