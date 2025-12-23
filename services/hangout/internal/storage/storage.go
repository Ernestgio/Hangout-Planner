package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, path string, reader io.Reader, contentType string) error

	Delete(ctx context.Context, path string) error

	GeneratePresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)
}
