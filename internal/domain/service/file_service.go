package service

import (
	"github.com/TancelinMazzotti/astigo/internal/domain/port/api/service"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/publisher"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/repository"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/spi/storage"
	"go.uber.org/zap"
)

var (
	_ service.IFileService = (*FileService)(nil)
)

type FileService struct {
	logger *zap.Logger
	repo   repository.IFileRepository
	store  storage.IFileStorage
	pub    publisher.IFilePublisher
}

func NewFileService(logger *zap.Logger, repo repository.IFileRepository, store storage.IFileStorage, pub publisher.IFilePublisher) *FileService {
	return &FileService{
		logger: logger,
		repo:   repo,
		store:  store,
		pub:    pub,
	}
}
