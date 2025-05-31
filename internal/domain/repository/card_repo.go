package repository

import (
	"context"

	"LostAndFound/internal/domain/entity"
)

type CardRepo interface {
	Create(ctx context.Context, l *entity.Card) error
	GetByID(ctx context.Context, id string) (*entity.Card, error)
	FindAll(ctx context.Context, filter string) ([]*entity.Card, error)
	Update(ctx context.Context, l *entity.Card) error
	Delete(ctx context.Context, id string) error
	FindNearLocation(ctx context.Context, lat, lon, radius float64, status string) ([]*entity.Card, error)
}
