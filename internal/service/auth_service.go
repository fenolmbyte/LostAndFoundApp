package service

import (
	"LostAndFound/internal/auth"
	"LostAndFound/internal/domain/entity"
	"LostAndFound/internal/domain/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     repository.UserRepo
	cacheRepo    repository.CacheRepo
	tokenManager *auth.TokenManager
}

func (a AuthService) Register(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	_, err := a.userRepo.FindByEmail(ctx, user.Email)
	if err == nil {
		return fmt.Errorf("this user already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt hashing failed: %w", err)
	}

	user.Password = string(hash)
	user.ID = uuid.New().String()

	return a.userRepo.Create(ctx, user)
}

func (a AuthService) Login(c context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "nil", fmt.Errorf("failed to find user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	role := "user"
	if user.IsAdmin {
		role = "admin"
	}

	token, err := a.tokenManager.Generate(user.ID, role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a AuthService) Logout(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()

	claims, err := a.tokenManager.Parse(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	if err = a.cacheRepo.BlacklistToken(ctx, token, ttl); err != nil {
		return fmt.Errorf("blacklist failed: %w", err)
	}

	if claims.UserID != "" {
		if err = a.cacheRepo.DeleteUserData(ctx, claims.UserID); err != nil {
			return fmt.Errorf("cache clear failed: %w", err)
		}
	}

	return a.cacheRepo.BlacklistToken(ctx, token, ttl)
}

func NewAuthService(userRepo repository.UserRepo, cacheRepo repository.CacheRepo, tokenManager *auth.TokenManager) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		cacheRepo:    cacheRepo,
		tokenManager: tokenManager,
	}
}
