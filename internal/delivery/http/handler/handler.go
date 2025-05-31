package handler

import (
	"LostAndFound/internal/auth"
	"LostAndFound/internal/service"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	services     *service.Service
	TokenManager *auth.TokenManager
	validator    *validator.Validate
}

func NewHandler(services *service.Service, tokenManager *auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		TokenManager: tokenManager,
		validator:    validator.New(),
	}
}
