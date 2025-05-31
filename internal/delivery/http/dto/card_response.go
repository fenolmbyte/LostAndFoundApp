package dto

import "time"

type OwnerDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Phone    string `json:"phone"`
	Telegram string `json:"telegram"`
}

type CardResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	DistanceM   float64   `json:"distance_m"`
	City        string    `json:"city"`
	Street      string    `json:"street"`
	PreviewURL  string    `json:"preview_url"`
	Images      []string  `json:"images"`
	Status      string    `json:"status"`
	Owner       OwnerDTO  `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
}
