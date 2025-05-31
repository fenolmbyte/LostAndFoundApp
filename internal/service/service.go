package service

import (
	"context"

	"LostAndFound/internal/auth"
	"LostAndFound/internal/bootstrap"
	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/domain/entity"
)

type Auth interface {
	Register(ctx context.Context, u *entity.User) error
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
}

type Users interface {
	GetProfile(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, u *entity.User) error
}

type Cards interface {
	CreateCard(ctx context.Context, l *entity.Card) error
	GetCardByID(ctx context.Context, id string) (*entity.Card, error)
	GetAllCards(ctx context.Context, filter string) ([]*entity.Card, error)
	UpdateCard(ctx context.Context, l *entity.Card) error
	DeleteCard(ctx context.Context, id string) error
	GetCardsNear(ctx context.Context, lat, lon, radius float64, status string) ([]*entity.Card, error)
}

type Files interface {
	GenerateUploadURL(ctx context.Context, userID string, req dto.FileRequest) (*dto.FileUploadResponse, error)
	DeleteFile(ctx context.Context, userID, key string) error
}

type Cache interface {
	SetUsername(ctx context.Context, userID, username string) error
	GetUsername(ctx context.Context, userID string) (string, error)
	DeleteUsername(ctx context.Context, userID string) error
}

type Service struct {
	Auth
	Users
	Cards
	Files
	Cache
}

func NewService(deps *bootstrap.Deps, tm *auth.TokenManager) *Service {
	return &Service{
		Auth:  NewAuthService(deps.UserRepo, deps.CacheRepo, tm),
		Users: NewUserService(deps.UserRepo),
		Cards: NewCardService(deps.CardRepo, deps.UserRepo, deps.CacheRepo, deps.FileStore),
		Files: NewFileService(deps.FileStore),
		Cache: NewCacheService(deps.CacheRepo),
	}
}
