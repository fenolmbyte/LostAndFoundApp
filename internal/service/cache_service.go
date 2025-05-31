package service

import (
	"LostAndFound/internal/domain/repository"
	"context"
)

type CacheService struct {
	repo repository.CacheRepo
}

func (c CacheService) SetUsername(ctx context.Context, userID string, username string) error {
	panic("implement me")
}

func (c CacheService) GetUsername(ctx context.Context, userID string) (string, error) {
	panic("implement me")
}

func (c CacheService) DeleteUsername(ctx context.Context, userID string) error {
	panic("implement me")
}

func NewCacheService(cacheRepo repository.CacheRepo) *CacheService {
	return &CacheService{repo: cacheRepo}
}
