package mapper

import (
	"time"

	"github.com/google/uuid"

	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/domain/entity"
)

func ToCardEntity(r dto.CreateCardRequest, ownerID string) *entity.Card {
	return &entity.Card{
		ID:          uuid.NewString(),
		Title:       r.Title,
		Description: r.Description,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		City:        r.City,
		Street:      r.Street,
		PreviewURL:  r.PreviewURL,
		Images:      r.Images,
		Status:      entity.CardStatus(r.Status),
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}
}

func ToCardUpdateEntity(dto dto.UpdateCardRequest, ownerID, cardID string) *entity.Card {
	return &entity.Card{
		ID:          cardID,
		Owner:       entity.Owner{ID: ownerID},
		Title:       dto.Title,
		Description: dto.Description,
		City:        dto.City,
		Street:      dto.Street,
		Status:      entity.CardStatus(dto.Status),
		PreviewURL:  dto.PreviewURL,
		Latitude:    dto.Latitude,
		Longitude:   dto.Longitude,
		Images:      dto.Images,
	}
}

func ToCardResponse(l *entity.Card, owner dto.OwnerDTO) dto.CardResponse {
	return dto.CardResponse{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		Latitude:    l.Latitude,
		Longitude:   l.Longitude,
		DistanceM:   l.DistanceM,
		City:        l.City,
		Street:      l.Street,
		PreviewURL:  l.PreviewURL,
		Images:      l.Images,
		Status:      string(l.Status),
		Owner:       owner,
		CreatedAt:   l.CreatedAt,
	}
}
