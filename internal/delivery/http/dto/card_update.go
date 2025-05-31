package dto

type UpdateCardRequest struct {
	Title       string   `json:"title,omitempty" validate:"omitempty,min=3"`
	Description string   `json:"description,omitempty" validate:"omitempty,min=10"`
	City        string   `json:"city,omitempty" validate:"omitempty"`
	Street      string   `json:"street,omitempty" validate:"omitempty"`
	Status      string   `json:"status,omitempty" validate:"omitempty,oneof=lost found"`
	PreviewURL  string   `json:"preview_url,omitempty" validate:"omitempty,url"`
	Latitude    float64  `json:"latitude,omitempty" validate:"omitempty"`
	Longitude   float64  `json:"longitude,omitempty" validate:"omitempty"`
	Images      []string `json:"images,omitempty" validate:"omitempty,dive,url"`
}
