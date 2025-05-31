package dto

type FileRequest struct {
	FileName    string `json:"file_name" validate:"required,min=1"`
	ContentType string `json:"content_type" validate:"required"`
}
