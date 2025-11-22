package model

import (
	"time"

	"github.com/google/uuid"
)

type UploadStatus string

const (
	UploadStatusPending  UploadStatus = "pending"
	UploadStatusComplete UploadStatus = "complete"
	UploadStatusFailed   UploadStatus = "failed"
)

type File struct {
	Id         uuid.UUID
	Name       string
	Size       int64
	Extension  string
	MimeType   string
	Path       string
	Status     UploadStatus
	CreatedAt  time.Time
	UploadedAt time.Time
	UpdatedAt  time.Time
}
