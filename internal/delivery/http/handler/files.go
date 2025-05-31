package handler

import (
	"LostAndFound/internal/delivery/http/dto"
	"encoding/json"
	"net/http"
)

// UploadFile godoc
// @Summary Генерация URL для загрузки файла
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.FileRequest true "Данные о файле"
// @Success 201 {object} dto.FileUploadResponse
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизован"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/files/upload [post]
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	var req dto.FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.validator.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	res, err := h.services.Files.GenerateUploadURL(r.Context(), userID, req)
	if err != nil {
		http.Error(w, "failed to generate URLs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {

}
