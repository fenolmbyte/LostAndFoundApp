package repository

import (
	"context"
	"time"

	"LostAndFound/internal/domain/entity"
)

type CacheRepo interface {
	SetUserData(ctx context.Context, user *entity.User) error
	GetUserData(ctx context.Context, userID string) (entity.User, error)
	DeleteUserData(ctx context.Context, userID string) error

	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)

	SaveCard(ctx context.Context, card *entity.Card) error
	GetCardByID(ctx context.Context, id string) (*entity.Card, error)
	DeleteCard(ctx context.Context, id string) error
}
