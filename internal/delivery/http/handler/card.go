package handler

import (
	e "LostAndFound/internal/common/errors"
	"LostAndFound/internal/delivery/http/dto"
	"LostAndFound/internal/delivery/http/mapper"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Summary Создание объявления
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateCardRequest true "Данные объявления"
// @Success 201 {string} string "Создано"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизован"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/cards [post]
func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.validator.Struct(req); err != nil {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	card := mapper.ToCardEntity(req, userID)

	if err := h.services.Cards.CreateCard(r.Context(), card); err != nil {
		http.Error(w, "failed to create card", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "card created successfully"})
}

// @Summary Получить объявление по ID
// @Tags Cards
// @Produce json
// @Param id path string true "ID объявления"
// @Success 200 {object} dto.CardResponse
// @Failure 404 {string} string "Объявление не найдено"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/cards/{id} [get]
func (h *Handler) GetCardByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	card, err := h.services.Cards.GetCardByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, e.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch card", http.StatusInternalServerError)
		return
	}
	if card == nil {
		http.Error(w, "card not found", http.StatusNotFound)
		return
	}

	owner := dto.OwnerDTO{
		ID:       card.Owner.ID,
		Name:     card.Owner.Name,
		Surname:  card.Owner.Surname,
		Phone:    card.Owner.Phone,
		Telegram: card.Owner.Telegram,
	}
	resp := mapper.ToCardResponse(card, owner)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Получить все объявления (по статусу)
// @Tags Cards
// @Produce json
// @Param status query string false "Статус объявления (lost/found)"
// @Success 200 {array} dto.CardResponse
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/cards/all [get]
func (h *Handler) GetAllCards(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")

	cards, err := h.services.GetAllCards(r.Context(), status)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cards: %v", err), http.StatusInternalServerError)
		return
	}

	var resp []dto.CardResponse
	for _, l := range cards {
		owner := dto.OwnerDTO{
			ID: l.Owner.ID, Name: l.Owner.Name,
			Surname: l.Owner.Surname, Phone: l.Owner.Phone,
			Telegram: l.Owner.Telegram,
		}
		resp = append(resp, mapper.ToCardResponse(l, owner))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Получить объявления поблизости
// @Tags Cards
// @Produce json
// @Param lat query number true "Широта"
// @Param lon query number true "Долгота"
// @Param radius query number true "Радиус поиска (км)"
// @Param status query string false "Статус"
// @Success 200 {array} dto.CardResponse
// @Failure 400 {string} string "Некорректные координаты"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/cards/near [get]
func (h *Handler) GetCardsNear(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	latStr := q.Get("lat")
	lonStr := q.Get("lon")
	radiusStr := q.Get("radius")
	status := q.Get("status")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		http.Error(w, "invalid radius", http.StatusBadRequest)
		return
	}

	cards, err := h.services.GetCardsNear(r.Context(), lat, lon, radius, status)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get nearby cards: %v", err), http.StatusInternalServerError)
		return
	}

	var resp []dto.CardResponse
	for _, l := range cards {
		owner := dto.OwnerDTO{
			ID: l.Owner.ID, Name: l.Owner.Name,
			Surname: l.Owner.Surname, Phone: l.Owner.Phone,
			Telegram: l.Owner.Telegram,
		}
		resp = append(resp, mapper.ToCardResponse(l, owner))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Обновить объявление
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID объявления"
// @Param input body dto.UpdateCardRequest true "Объявление"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизован"
// @Failure 403 {string} string "Нет доступа"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/cards/{id} [put]
func (h *Handler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var card dto.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		http.Error(w, "missing card ID", http.StatusBadRequest)
		return
	}

	entity := mapper.ToCardUpdateEntity(card, userId, cardID)
	entity.ID = cardID

	if err := h.services.Cards.UpdateCard(r.Context(), entity); err != nil {
		http.Error(w, "failed to update card: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "card updated successfully"})
}

// @Summary Удалить объявление
// @Tags Cards
// @Produce json
// @Param id path string true "ID объявления"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 403 {string} string "Нет доступа"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/cards/{id} [delete]
func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	cardID := chi.URLParam(r, "id")

	if cardID == "" {
		http.Error(w, "missing card ID", http.StatusBadRequest)
		return
	}

	if err := h.services.Cards.DeleteCard(r.Context(), cardID); err != nil {
		if errors.Is(err, e.ErrPermissionDenied) {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		http.Error(w, "failed to delete card: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "card deleted successfully"})
}
