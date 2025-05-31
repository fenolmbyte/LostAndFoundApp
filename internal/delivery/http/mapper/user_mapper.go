package mapper

import (
	"time"

	"github.com/google/uuid"

	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/domain/entity"
)

func ToUserEntity(r dto.UserRegisterRequest) *entity.User {
	return &entity.User{
		ID:        uuid.New().String(),
		Email:     r.Email,
		Password:  r.Password,
		Name:      r.Name,
		Surname:   r.Surname,
		Phone:     r.Phone,
		Telegram:  r.Telegram,
		CreatedAt: time.Now(),
	}
}

func ToUserDTO(u *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Name:     u.Name,
		Surname:  u.Surname,
		Phone:    u.Phone,
		Telegram: u.Telegram,
	}
}
