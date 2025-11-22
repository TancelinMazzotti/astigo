package service

import (
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/service"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/messaging"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/repository"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/storage"
	"go.uber.org/zap"
)

var (
	_ service.IFileService = (*FileService)(nil)
)

type FileService struct {
	logger *zap.Logger
	repo   repository.IFileRepository
	store  storage.IFileStorage
	msg    messaging.IFileMessaging
}

func NewFileService(logger *zap.Logger, repo repository.IFileRepository, store storage.IFileStorage, msg messaging.IFileMessaging) *FileService {
	return &FileService{
		logger: logger,
		repo:   repo,
		store:  store,
		msg:    msg,
	}
}
