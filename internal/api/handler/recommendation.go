package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/service"
)

type RecommendationHandler struct {
	BaseHandler
	recommendationService service.RecommendationService
}

func NewRecommendationHandler(recommendationService service.RecommendationService, logger zerolog.Logger) *RecommendationHandler {
	return &RecommendationHandler{
		BaseHandler:           NewBaseHandler(logger),
		recommendationService: recommendationService,
	}
}

func (h *RecommendationHandler) RegisterRoutes(r chi.Router) {
	r.Get("/daily/{profileId}", h.GetDailyRecommendations)
	r.Get("/meal-plan/{profileId}", h.GetMealPlanRecommendations)
	r.Get("/alternatives/{foodId}", h.GetFoodAlternatives)
}

func (h *RecommendationHandler) GetDailyRecommendations(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileId")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 5
	}

	recommendations, err := h.recommendationService.GetDailyRecommendations(r.Context(), profileID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"recommendations": recommendations,
		"profileId":       profileID,
		"limit":           limit,
	})
}

func (h *RecommendationHandler) GetMealPlanRecommendations(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileId")
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 7
	}

	mealPlan, err := h.recommendationService.GetMealPlanRecommendations(r.Context(), profileID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"mealPlan":  mealPlan,
		"profileId": profileID,
		"days":      days,
	})
}

func (h *RecommendationHandler) GetFoodAlternatives(w http.ResponseWriter, r *http.Request) {
	foodID := chi.URLParam(r, "foodId")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 5
	}

	alternatives, err := h.recommendationService.GetFoodAlternatives(r.Context(), foodID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"alternatives": alternatives,
		"foodId":       foodID,
		"limit":        limit,
	})
}
