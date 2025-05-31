package service

import (
	e "LostAndFound/internal/common/errors"
	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/domain/repository"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	repo repository.FileStorage
}

func (f *FileService) GenerateUploadURL(ctx context.Context, userID string, req dto.FileRequest) (*dto.FileUploadResponse, error) {
	if strings.Contains(req.FileName, "..") || strings.Contains(req.FileName, "/") {
		return nil, fmt.Errorf("invalid file name")
	}

	ext := filepath.Ext(req.FileName)
	key := fmt.Sprintf("users/%s/%s%s", userID, uuid.New().String(), ext)

	presignedURL, err := f.repo.GeneratePresignedPutURL(key, req.ContentType, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s/%s", f.repo.GetBaseURL(), f.repo.GetBucket(), key)

	return &dto.FileUploadResponse{
		FileName:     req.FileName,
		PresignedURL: presignedURL,
		PublicURL:    publicURL,
	}, nil
}

func (f *FileService) DeleteFile(ctx context.Context, userID string, key string) error {
	if !strings.HasPrefix(key, "users/"+userID+"/") {
		return fmt.Errorf("forbidden")
	}

	exists, err := f.repo.FileExists(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}
	if !exists {
		return e.ErrNotFound
	}

	if err = f.repo.DeleteFile(ctx, key); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func NewFileService(fileRepo repository.FileStorage) *FileService {
	return &FileService{repo: fileRepo}
}
