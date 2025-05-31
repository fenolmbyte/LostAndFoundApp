package service

import (
	e "LostAndFound/internal/common/errors"
	"LostAndFound/internal/domain/entity"
	"LostAndFound/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepo
}

func (u *UserService) GetProfile(c context.Context, userID string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot find user")
	}
	return &entity.User{
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		Phone:     user.Phone,
		Telegram:  user.Telegram,
		CreatedAt: user.CreatedAt,
		IsAdmin:   user.IsAdmin,
	}, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, updated *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	currentUser, err := u.repo.FindByID(ctx, updated.ID)
	if err != nil {
		return e.ErrNotFound
	}

	changed := false

	if len(updated.Email) != 0 && currentUser.Email != updated.Email {
		currentUser.Email = updated.Email
		changed = true
	}
	if len(updated.Name) != 0 && currentUser.Name != updated.Name {
		currentUser.Name = updated.Name
		changed = true
	}
	if len(updated.Surname) != 0 && currentUser.Surname != updated.Surname {
		currentUser.Surname = updated.Surname
		changed = true
	}
	if len(updated.Phone) != 0 && currentUser.Phone != updated.Phone {
		currentUser.Phone = updated.Phone
		changed = true
	}
	if len(updated.Telegram) != 0 && currentUser.Telegram != updated.Telegram {
		currentUser.Telegram = updated.Telegram
		changed = true
	}
	if updated.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updated.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("bcrypt hashing failed: %w", err)
		}
		currentUser.Password = string(hashedPassword)
		changed = true
	}

	if !changed {
		return e.ErrNoChanges
	}

	return u.repo.Update(ctx, currentUser)
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{repo: userRepo}
}
