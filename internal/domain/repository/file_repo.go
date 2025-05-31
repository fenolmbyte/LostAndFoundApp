package repository

import (
	"context"
	"time"
)

type FileStorage interface {
	GeneratePresignedPutURL(key, contentType string, expires time.Duration) (string, error)
	DeleteFile(ctx context.Context, key string) error
	FileExists(ctx context.Context, key string) (bool, error)
	GetBaseURL() string
	GetBucket() string
}
