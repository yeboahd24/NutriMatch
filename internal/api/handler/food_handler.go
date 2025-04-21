package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api/middleware/auth"
	"github.com/yeboahd24/nutrimatch/internal/config"
	"github.com/yeboahd24/nutrimatch/internal/service"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
	"github.com/yeboahd24/nutrimatch/pkg/response"
)

type FoodHandler struct {
	BaseHandler
	foodService service.FoodService
	validator   *validator.Validate
	jwtConfig   config.JWTConfig
}

func NewFoodHandler(foodService service.FoodService, logger zerolog.Logger, jwtConfig config.JWTConfig) *FoodHandler {
	return &FoodHandler{
		BaseHandler: NewBaseHandler(logger),
		foodService: foodService,
		validator:   validator.New(),
		jwtConfig:   jwtConfig,
	}
}

func (h *FoodHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.SearchFoods)
	r.Get("/{id}", h.GetFood)
	r.Get("/category/{category}", h.GetFoodsByCategory)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(h.jwtConfig))
		r.Post("/{foodId}/rate", h.RateFood)
		r.Get("/ratings", h.ListUserRatings)
		r.Get("/saved", h.ListSavedFoods)
		r.Post("/{foodId}/save", h.SaveFood)
		r.Delete("/{foodId}/save", h.RemoveSavedFood)
	})
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
		h.logger.Error().Err(err).Str("query", query).Msg("Failed to search foods")
		response.Error(w, apperrors.Internal("Failed to search foods", err))
		return
	}

	meta := response.PaginationMeta(page, limit, total)
	response.JSONWithMeta(w, http.StatusOK, foods, meta)
}

func (h *FoodHandler) GetFood(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	food, err := h.foodService.GetFood(r.Context(), id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.Error(w, apperrors.NotFound("food", err))
			return
		}
		h.logger.Error().Err(err).Str("food_id", id).Msg("Failed to get food")
		response.Error(w, apperrors.Internal("Failed to get food", err))
		return
	}

	response.JSON(w, http.StatusOK, food)
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
		h.logger.Error().Err(err).Str("category", category).Msg("Failed to get foods by category")
		response.Error(w, apperrors.Internal("Failed to get foods by category", err))
		return
	}

	meta := response.PaginationMeta(page, limit, total)
	data := map[string]interface{}{
		"foods":    foods,
		"category": category,
	}

	response.JSONWithMeta(w, http.StatusOK, data, meta)
}

func (h *FoodHandler) RateFood(w http.ResponseWriter, r *http.Request) {
	foodID := chi.URLParam(r, "foodId")
	if foodID == "" {
		response.Error(w, apperrors.InvalidInput("Food ID is required", nil))
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		response.Error(w, apperrors.Unauthorized("Unauthorized", nil))
		return
	}

	var input struct {
		Rating   int    `json:"rating" validate:"required,min=1,max=5"`
		Comments string `json:"comments,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request payload", err))
		return
	}

	if err := h.validator.Struct(input); err != nil {
		response.Error(w, apperrors.InvalidInput("Validation failed", err))
		return
	}

	rating, err := h.foodService.RateFood(r.Context(), userID, foodID, input.Rating, input.Comments)
	if err != nil {
		if err.Error() == "food not found" {
			response.Error(w, apperrors.NotFound("Food not found", err))
			return
		}
		if err.Error() == "user has already rated this food" {
			response.Error(w, apperrors.DuplicateEntity("User has already rated this food", err))
			return
		}
		h.logger.Error().Err(err).
			Str("user_id", userID.String()).
			Str("food_id", foodID).
			Int("rating", input.Rating).
			Msg("Failed to rate food")
		response.Error(w, apperrors.Internal("Failed to rate food", err))
		return
	}

	response.JSON(w, http.StatusOK, rating)
}

func (h *FoodHandler) ListUserRatings(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		response.Error(w, apperrors.Unauthorized("Unauthorized", nil))
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 {
		limit = 10
	}

	ratings, err := h.foodService.ListUserRatings(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to list user ratings")
		response.Error(w, apperrors.Internal("Failed to list ratings", err))
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"ratings": ratings,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
			},
		},
	})
}

func (h *FoodHandler) SaveFood(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		response.Error(w, apperrors.Unauthorized("Unauthorized", nil))
		return
	}

	foodID := chi.URLParam(r, "foodId")
	if foodID == "" {
		response.Error(w, apperrors.InvalidInput("Food ID is required", nil))
		return
	}

	var input struct {
		ListType string `json:"list_type" validate:"required,oneof=favorites shopping_list watch_list"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request body", err))
		return
	}

	savedFood, err := h.foodService.SaveFood(r.Context(), userID, foodID, input.ListType)
	if err != nil {
		h.logger.Error().Err(err).
			Str("user_id", userID.String()).
			Str("food_id", foodID).
			Str("list_type", input.ListType).
			Msg("Failed to save food")
		response.Error(w, apperrors.Internal("Failed to save food", err))
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    savedFood,
	})
}

func (h *FoodHandler) ListSavedFoods(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		response.Error(w, apperrors.Unauthorized("Unauthorized", nil))
		return
	}

	listType := r.URL.Query().Get("list_type")
	if listType == "" {
		listType = "favorites" // Default list type
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 {
		limit = 10
	}

	savedFoods, err := h.foodService.ListSavedFoods(r.Context(), userID, listType, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("user_id", userID.String()).
			Str("list_type", listType).
			Msg("Failed to list saved foods")
		response.Error(w, apperrors.Internal("Failed to list saved foods", err))
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"saved_foods": savedFoods,
			"list_type":   listType,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
			},
		},
	})
}

func (h *FoodHandler) RemoveSavedFood(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		response.Error(w, apperrors.Unauthorized("Unauthorized", nil))
		return
	}

	foodID := chi.URLParam(r, "foodId")
	if foodID == "" {
		response.Error(w, apperrors.InvalidInput("Food ID is required", nil))
		return
	}

	listType := r.URL.Query().Get("list_type")
	if listType == "" {
		listType = "favorites" // Default list type
	}

	if err := h.foodService.RemoveSavedFood(r.Context(), userID, foodID, listType); err != nil {
		h.logger.Error().Err(err).
			Str("user_id", userID.String()).
			Str("food_id", foodID).
			Str("list_type", listType).
			Msg("Failed to remove saved food")
		response.Error(w, apperrors.Internal("Failed to remove saved food", err))
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Food removed from saved list",
	})
}
