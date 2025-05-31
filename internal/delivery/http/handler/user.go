package handler

import (
	e "LostAndFound/internal/common/errors"
	utils "LostAndFound/internal/common/validation"
	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/delivery/http/mapper"
	"encoding/json"
	"errors"
	"net/http"
)

// GetProfile получает профиль текущего пользователя
// @Summary Получить свой профиль
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /users [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(string)
	if !ok || id == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.services.Users.GetProfile(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	userDTO := mapper.ToUserDTO(user)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDTO)
}

// GetProfileByID получает профиль пользователя по ID
// @Summary Получить профиль пользователя по ID
// @Tags users
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {string} string "User not found"
// @Router /users/profile [get]
func (h *Handler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")

	user, err := h.services.Users.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	userDTO := mapper.ToUserDTO(user)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDTO)
}

// UpdateProfile обновляет профиль текущего пользователя
// @Summary Обновить свой профиль
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body dto.UpdateUserRequest true "Данные для обновления"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid request or no changes"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error updating user"
// @Router /users/update [put]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(string)
	if !ok || id == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, utils.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	user := mapper.ToUserEntity(dto.UserRegisterRequest(req))
	user.ID = id

	if err := h.services.Users.UpdateProfile(r.Context(), user); err != nil {
		if errors.Is(err, e.ErrNoChanges) {
			http.Error(w, "no changes made to profile", http.StatusBadRequest)
			return
		}
		if errors.Is(err, e.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "profile updated successfully"})
}
