package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, path string, reader io.Reader, contentType string) error
	Delete(ctx context.Context, path string) error
	GeneratePresignedDownloadURL(ctx context.Context, path string) (string, error)
	GeneratePresignedUploadURL(ctx context.Context, path string, contentType string) (string, error)
}
