package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api/middleware/auth"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
	"github.com/yeboahd24/nutrimatch/internal/service"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
	"github.com/yeboahd24/nutrimatch/pkg/response"
)

type ProfileHandler struct {
	BaseHandler
	profileService service.ProfileService
	validator      *validator.Validate
}

func NewProfileHandler(profileService service.ProfileService, logger zerolog.Logger) *ProfileHandler {
	return &ProfileHandler{
		BaseHandler:    NewBaseHandler(logger),
		profileService: profileService,
		validator:      validator.New(),
	}
}

func (h *ProfileHandler) RegisterRoutes(r chi.Router) {
	// Specific routes first
	r.Get("/user/me", h.ListUserProfiles)      // List all profiles for the authenticated user
	r.Get("/debug/all", h.GetAllProfiles)      // Debug endpoint to get all profiles
	r.Get("/check/{id}", h.CheckProfileExists) // Debug endpoint

	// Wildcard routes last
	r.Post("/{userId}", h.CreateProfile)
	r.Get("/{id}", h.GetProfile)
	r.Put("/{id}", h.UpdateProfile)
	r.Delete("/{id}", h.DeleteProfile)
}

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL parameter
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		appErr := apperrors.InvalidInput("User ID is required", nil)
		appErr.WriteJSON(w)
		return
	}

	var req struct {
		Age         int      `json:"age"`
		Gender      string   `json:"gender"`
		Weight      float64  `json:"weight"`
		Height      float64  `json:"height"`
		Goals       []string `json:"goals"`
		Allergies   []string `json:"allergies"`
		Preferences []string `json:"preferences"`
		IsDefault   bool     `json:"is_default"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appErr := apperrors.InvalidInput("Invalid request body", err)
		appErr.WriteJSON(w)
		return
	}

	// Log the request for debugging
	h.logger.Debug().Str("user_id", userID).Interface("goals", req.Goals).Interface("allergies", req.Allergies).Bool("is_default", req.IsDefault).Msg("Creating profile with user ID from URL parameter")

	// Use the user ID from the URL parameter
	profile, err := h.profileService.CreateProfile(r.Context(), userID, req.Age, req.Gender, req.Weight, req.Height, req.Goals, req.Allergies, req.Preferences, req.IsDefault)
	if err != nil {
		// Handle database errors
		if err.Error() == "sql: no rows in result set" {
			appErr := apperrors.NotFound("User not found", err)
			appErr.WriteJSON(w)
			return
		}

		appErr := apperrors.Internal("Failed to create profile", err)
		appErr.WriteJSON(w)
		return
	}

	// The profile ID has been updated with the database-generated ID in the repository layer
	// Log the final profile ID for debugging
	h.logger.Info().Str("profile_id", profile.ID.String()).Msg("Profile created with database-generated ID")

	// Verify the profile was actually created by retrieving it
	verifiedProfile, err := h.profileService.GetProfile(r.Context(), profile.ID.String())
	if err != nil {
		h.logger.Error().Err(err).Str("profile_id", profile.ID.String()).Msg("Failed to verify profile creation")
		// Continue anyway, but log the error
	} else {
		h.logger.Info().Str("profile_id", profile.ID.String()).Str("verified_id", verifiedProfile.ID.String()).Msg("Profile creation verified")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.Error(w, apperrors.NotFound("profile not found", err))
			return
		}
		h.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get profile")
		response.Error(w, apperrors.Internal("Failed to get profile", err))
		return
	}

	response.JSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var input profile.Profile
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request payload", err))
		return
	}

	userID := r.Context().Value("user_id").(string)
	input.UserID = userID

	if err := h.validator.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var details []map[string]string
			for _, err := range validationErrors {
				details = append(details, map[string]string{
					"field":   err.Field(),
					"message": err.Tag(),
				})
			}
			response.Error(w, apperrors.InvalidInputWithDetails("Validation failed", details, err))
			return
		}
		response.Error(w, apperrors.InvalidInput("Validation failed", err))
		return
	}

	err := h.profileService.UpdateProfile(r.Context(), userID, input.Age, input.Gender, input.Weight, input.Height, input.Goals, input.Allergies, input.Preferences, input.IsDefault)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to update profile")
		response.Error(w, apperrors.Internal("Failed to update profile", err))
		return
	}

	updatedProfile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, apperrors.Internal("Failed to retrieve updated profile", err))
		return
	}

	response.JSON(w, http.StatusOK, updatedProfile)
}

func (h *ProfileHandler) CheckProfileExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info().Str("profile_id", id).Msg("Checking if profile exists")

	// Try to get the profile from the service
	profile, err := h.profileService.GetProfile(r.Context(), id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// Try to get all profiles for debugging
			allProfiles, _ := h.profileService.GetAllProfiles(r.Context())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"exists":       false,
				"error":        "Profile not found",
				"all_profiles": allProfiles,
				"debug_profile": map[string]interface{}{
					"id":                   "f6335a0d-2835-4c81-bb74-25a568429f64",
					"user_id":              "18405c62-70b4-44ef-b84a-22a076130b57",
					"profile_name":         "",
					"is_default":           false,
					"health_conditions":    []string{},
					"dietary_restrictions": []string{},
					"allergens":            []string{"peanuts", "shellfish"},
					"goal_type":            "weight_loss",
					"disliked_foods":       []string{},
					"preferred_foods":      []string{"vegetarian", "low_carb"},
					"cuisine_preferences":  []string{},
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists": false,
			"error":  err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exists":  true,
		"profile": profile,
	})
}

// ListUserProfiles returns all profiles for the authenticated user
func (h *ProfileHandler) ListUserProfiles(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the context
	userID, ok := auth.GetUserID(r)
	if !ok {
		appErr := apperrors.Unauthorized("User ID not found in context", nil)
		appErr.WriteJSON(w)
		return
	}

	h.logger.Info().Str("user_id", userID.String()).Msg("Getting all profiles for authenticated user")

	// Get all profiles for the user
	profiles, err := h.profileService.GetProfilesByUserID(r.Context(), userID.String())
	if err != nil {
		// Handle database errors
		if err.Error() == "sql: no rows in result set" {
			// No profiles found is not an error, just return an empty array
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"count":    0,
				"profiles": []interface{}{},
			})
			return
		}

		h.logger.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get profiles for user")
		appErr := apperrors.Internal("Failed to get user profiles", err)
		appErr.WriteJSON(w)
		return
	}

	h.logger.Info().Str("user_id", userID.String()).Int("profile_count", len(profiles)).Msg("Successfully retrieved profiles for user")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":    len(profiles),
		"profiles": profiles,
	})
}

// GetAllProfiles is a debug endpoint to get all profiles in the database
func (h *ProfileHandler) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Getting all profiles for debugging")

	// Get all profiles from the database
	profiles, err := h.profileService.GetAllProfiles(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get all profiles")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":    len(profiles),
		"profiles": profiles,
	})
}

func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get the authenticated user ID from the context
	userID, ok := auth.GetUserID(r)
	if !ok {
		appErr := apperrors.Unauthorized("User ID not found in context", nil)
		appErr.WriteJSON(w)
		return
	}

	// First, check if the profile exists and belongs to the authenticated user
	profile, err := h.profileService.GetProfile(r.Context(), id)
	if err != nil {
		// Handle database errors
		if err.Error() == "sql: no rows in result set" {
			appErr := apperrors.NotFound("Profile not found", err)
			appErr.WriteJSON(w)
			return
		}

		appErr := apperrors.Internal("Failed to get profile", err)
		appErr.WriteJSON(w)
		return
	}

	// Verify ownership
	if profile.UserID.String() != userID.String() {
		appErr := apperrors.Forbidden("You don't have permission to delete this profile", nil)
		appErr.WriteJSON(w)
		return
	}

	err = h.profileService.DeleteProfile(r.Context(), id)
	if err != nil {
		// Handle database errors
		if err.Error() == "sql: no rows in result set" {
			appErr := apperrors.NotFound("Profile not found", err)
			appErr.WriteJSON(w)
			return
		}

		appErr := apperrors.Internal("Failed to delete profile", err)
		appErr.WriteJSON(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
