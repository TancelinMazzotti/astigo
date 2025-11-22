package storage

import (
	"context"
	"net/http"
)

type IFileStorage interface {
	PresignedDownloadURL(ctx context.Context, path string) (string, http.Header, error)
	PresignedUploadURL(ctx context.Context, key, contentType string) (string, http.Header, error)
	Delete(ctx context.Context, path string) error
}
