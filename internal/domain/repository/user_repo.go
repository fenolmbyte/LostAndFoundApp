package repository

import (
	"context"

	"LostAndFound/internal/domain/entity"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	Delete(ctx context.Context, id string) error
}
