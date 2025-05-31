package dto

type FileUploadResponse struct {
	FileName     string `json:"file_name"`
	PresignedURL string `json:"presigned_url"`
	PublicURL    string `json:"public_url"`
}
