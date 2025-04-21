package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api/middleware/auth"
	"github.com/yeboahd24/nutrimatch/internal/domain/recommendation"
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
	r.Get("/", h.GetRecommendations)
	r.Post("/filter", h.FilterRecommendations)
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
		h.logger.Error().Err(err).Msg("Failed to get daily recommendations")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"recommendations": recommendations,
			"profile_id":      profileID,
			"limit":           limit,
		},
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
		h.logger.Error().Err(err).Msg("Failed to get meal plan recommendations")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"meal_plan":  mealPlan,
			"profile_id": profileID,
			"days":       days,
		},
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
		h.logger.Error().Err(err).Msg("Failed to get food alternatives")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"alternatives": alternatives,
			"food_id":      foodID,
			"limit":        limit,
		},
	})
}

func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context using the auth middleware helper
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	profileID := r.URL.Query().Get("profileId")

	var profileUUID *uuid.UUID
	if profileID != "" {
		parsed, err := uuid.Parse(profileID)
		if err != nil {
			http.Error(w, "invalid profile ID", http.StatusBadRequest)
			return
		}
		profileUUID = &parsed
	}

	// Create recommendation request
	req := recommendation.RecommendationRequest{
		ProfileID: profileUUID,
		Limit:     limit,
		Offset:    offset,
	}

	// Get recommendations
	resp, err := h.recommendationService.GetRecommendations(userID, req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recommendations")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"recommendations": resp.Foods,
			"total_count":     resp.TotalCount,
			"applied_rules":   resp.AppliedRules,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
			},
		},
	})
}

func (h *RecommendationHandler) FilterRecommendations(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context using the auth middleware helper
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req recommendation.RecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get filtered recommendations
	resp, err := h.recommendationService.GetRecommendations(userID, req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to filter recommendations")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"recommendations": resp.Foods,
			"total_count":     resp.TotalCount,
			"applied_rules":   resp.AppliedRules,
			"pagination": map[string]interface{}{
				"limit":  req.Limit,
				"offset": req.Offset,
			},
		},
	})
}
