package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/service"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
)

type FoodHandler struct {
	BaseHandler
	foodService service.FoodService
}

func NewFoodHandler(foodService service.FoodService, logger zerolog.Logger) *FoodHandler {
	return &FoodHandler{
		BaseHandler: NewBaseHandler(logger),
		foodService: foodService,
	}
}

func (h *FoodHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.SearchFoods)
	r.Get("/{id}", h.GetFood)
	r.Get("/category/{category}", h.GetFoodsByCategory)
}

func (h *FoodHandler) SearchFoods(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	foods, total, err := h.foodService.SearchFoods(r.Context(), query, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"foods": foods,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *FoodHandler) GetFood(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	food, err := h.foodService.GetFood(r.Context(), id)
	if err != nil {
		// Handle database errors
		if err.Error() == "sql: no rows in result set" {
			appErr := apperrors.NotFound("Food not found", err)
			appErr.WriteJSON(w)
			return
		}

		appErr := apperrors.Internal("Failed to get food", err)
		appErr.WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(food)
}

func (h *FoodHandler) GetFoodsByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	foods, total, err := h.foodService.GetFoodsByCategory(r.Context(), category, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"foods":    foods,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"category": category,
	})
}
