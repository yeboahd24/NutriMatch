package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api/middleware/auth"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
	"github.com/yeboahd24/nutrimatch/internal/service"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
)

type UserHandler struct {
	BaseHandler
	userService service.UserService
}

func NewUserHandler(userService service.UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(logger),
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get("/me", h.GetProfile)
	r.Put("/me", h.UpdateProfile)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context using the auth middleware helper
	userID, ok := auth.GetUserID(r)
	if !ok {
		ErrorResponse(w, r, apperrors.Unauthorized("User ID not found in context", nil), h.logger)
		return
	}

	// Convert UUID to string for the service call
	profile, err := h.userService.GetProfile(r.Context(), userID.String())
	if err != nil {
		// Handle database errors specifically
		if err.Error() == "sql: no rows in result set" {
			HandleDatabaseError(w, r, err, "User profile", h.logger)
			return
		}

		ErrorResponse(w, r, apperrors.Internal("Failed to get user profile", err), h.logger)
		return
	}

	RespondOK(w, profile)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context using the auth middleware helper
	userID, ok := auth.GetUserID(r)
	if !ok {
		ErrorResponse(w, r, apperrors.Unauthorized("User ID not found in context", nil), h.logger)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := DecodeJSONBody(w, r, &req); err != nil {
		ErrorResponse(w, r, err, h.logger)
		return
	}

	// Log the update attempt
	h.logger.Info().Str("user_id", userID.String()).Str("name", req.Name).Str("email", req.Email).Msg("User profile update requested")

	err := h.userService.UpdateProfile(r.Context(), userID.String(), req.Name, req.Email)
	if err != nil {
		// Handle email update not allowed error
		if err == user.ErrEmailUpdateNotAllowed {
			ErrorResponse(w, r, apperrors.InvalidInput("Email updates are not allowed for security reasons. Please contact support if you need to change your email.", err), h.logger)
			return
		}

		// Handle database errors specifically
		if err.Error() == "sql: no rows in result set" {
			HandleDatabaseError(w, r, err, "User profile", h.logger)
			return
		}

		ErrorResponse(w, r, apperrors.Internal("Failed to update user profile", err), h.logger)
		return
	}

	h.logger.Info().Str("user_id", userID.String()).Msg("User profile updated successfully")

	RespondNoContent(w)
}
