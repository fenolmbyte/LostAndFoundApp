package dto

type CreateCardRequest struct {
	Title       string   `json:"title"       validate:"required"`
	Description string   `json:"description"`
	Latitude    float64  `json:"latitude"    validate:"required"`
	Longitude   float64  `json:"longitude"   validate:"required"`
	PreviewURL  string   `json:"preview_url"`
	Status      string   `json:"status"      validate:"required,oneof=lost found"`
	Images      []string `json:"images"`
	City        string   `json:"city"        validate:"required"`
	Street      string   `json:"street"`
}
