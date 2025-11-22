package repository

import (
	"context"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/data"
	"github.com/google/uuid"
)

type IFileRepository interface {
	FindAll(ctx context.Context, pagination data.PaginationOffset) ([]*model.File, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.File, error)
	Create(ctx context.Context, file *model.File) error
	Update(ctx context.Context, file *model.File) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
