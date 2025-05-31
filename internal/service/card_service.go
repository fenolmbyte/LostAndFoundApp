package service

import (
	e "LostAndFound/internal/common/errors"
	"LostAndFound/internal/domain/entity"
	"LostAndFound/internal/domain/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type CardService struct {
	userRepo  repository.UserRepo
	repo      repository.CardRepo
	cacheRepo repository.CacheRepo
	fileRepo  repository.FileStorage
}

func (l *CardService) CreateCard(c context.Context, card *entity.Card) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	userID, ok := c.Value("userID").(string)
	if !ok || userID == "" {
		return e.ErrUnauthorized
	}
	card.Owner.ID = userID

	owner, err := l.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error finding owner by id: %w", err)
	}
	card.Owner.Name = owner.Name
	card.Owner.Surname = owner.Surname
	card.Owner.Phone = owner.Phone
	card.Owner.Telegram = owner.Telegram
	card.ID = uuid.New().String()

	return l.repo.Create(ctx, card)
}

func (l *CardService) GetCardByID(c context.Context, id string) (*entity.Card, error) {
	ctx, cancel := context.WithTimeout(c, 50*time.Second)
	defer cancel()

	card, err := l.cacheRepo.GetCardByID(ctx, id)
	if err == nil && card != nil {
		return card, nil
	}

	card, err = l.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get card from DB: %w", err)
	}

	_ = l.cacheRepo.SaveCard(ctx, card)

	return card, nil
}

func (l *CardService) GetAllCards(ctx context.Context, filter string) ([]*entity.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return l.repo.FindAll(ctx, filter)
}

func (l *CardService) UpdateCard(c context.Context, updated *entity.Card) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	current, err := l.repo.GetByID(ctx, updated.ID)
	if err != nil {
		return e.ErrNotFound
	}

	changed := false

	if updated.Title != "" && updated.Title != current.Title {
		current.Title = updated.Title
		changed = true
	}
	if updated.Description != "" && updated.Description != current.Description {
		current.Description = updated.Description
		changed = true
	}
	if updated.City != "" && updated.City != current.City {
		current.City = updated.City
		changed = true
	}
	if updated.Street != "" && updated.Street != current.Street {
		current.Street = updated.Street
		changed = true
	}
	if updated.Status != "" && updated.Status != current.Status {
		current.Status = updated.Status
		changed = true
	}
	if updated.PreviewURL != "" && updated.PreviewURL != current.PreviewURL {
		current.PreviewURL = updated.PreviewURL
		changed = true
	}
	if updated.Latitude != 0 && updated.Latitude != current.Latitude {
		current.Latitude = updated.Latitude
		changed = true
	}
	if updated.Longitude != 0 && updated.Longitude != current.Longitude {
		current.Longitude = updated.Longitude
		changed = true
	}
	if len(updated.Images) > 0 && !slices.Equal(updated.Images, current.Images) {
		current.Images = updated.Images
		changed = true
	}

	if !changed {
		return e.ErrNoChanges
	}

	if err = l.repo.Update(ctx, current); err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	_ = l.cacheRepo.DeleteCard(ctx, current.ID)

	return nil
}

func (l *CardService) DeleteCard(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	userID := ctx.Value("userID").(string)

	card, err := l.GetCardByID(ctx, id)
	if err != nil {
		return err
	}
	if userID != card.Owner.ID {
		return e.ErrPermissionDenied
	}
	if err = l.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	_ = l.cacheRepo.DeleteCard(ctx, id)

	for _, key := range card.Images {
		if err := l.fileRepo.DeleteFile(ctx, key); err != nil {

			fmt.Printf("failed to delete file %s: %v\n", key, err)
		}
	}

	return nil
}

func (l *CardService) GetCardsNear(ctx context.Context, lat, lon, radius float64, status string) ([]*entity.Card, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return l.repo.FindNearLocation(ctx, lat, lon, radius, status)
}

func NewCardService(cardRepo repository.CardRepo, userRepo repository.UserRepo, cache repository.CacheRepo, fileRepo repository.FileStorage) *CardService {
	return &CardService{
		repo:      cardRepo,
		userRepo:  userRepo,
		cacheRepo: cache,
		fileRepo:  fileRepo,
	}
}
