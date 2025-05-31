package handler

import (
	v "LostAndFound/internal/common/validation"
	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/delivery/http/mapper"
	"encoding/json"
	"net/http"
)

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Регистрирует нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.UserRegisterRequest  true  "Данные пользователя"
// @Success      201    {string}  string                       "created"
// @Failure      400    {string}  string                       "invalid request"
// @Failure      500    {string}  string                       "internal error"
// @Router       /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.validator.Struct(req); err != nil {
		http.Error(w, v.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	user := mapper.ToUserEntity(req)

	if err := h.services.Auth.Register(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "you are registered successfully"})
}

// Login godoc
// @Summary      Логин пользователя
// @Description  Авторизует пользователя и возвращает токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.UserAuthRequest  true  "Email и пароль"
// @Success      200    {object}  string           "JWT токены"
// @Failure      400    {string}  string
// @Failure      401    {string}  string  "Invalid credentials"
// @Router       /auth/login [post]

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, v.FormatValidationError(err), http.StatusBadRequest)
	}

	token, err := h.services.Auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "you are logged in successfully", "token": token})
}

// Logout godoc
// @Summary      Выход из аккаунта
// @Description  Инвалидирует JWT токен
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := h.TokenManager.GetToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err = h.services.Auth.Logout(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "you are logged out successfully"})
}
