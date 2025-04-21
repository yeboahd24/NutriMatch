package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/service"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	BaseHandler
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(logger),
		authService: authService,
	}
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/refresh", h.RefreshToken)
	r.Post("/logout", h.Logout)
}

// @Summary Login user
// @Description Authenticate a user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body docs.LoginRequest true "User credentials"
// @Success 200 {object} docs.Response{data=docs.TokenResponse}
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := DecodeJSONBody(w, r, &req); err != nil {
		ErrorResponse(w, r, err, h.logger)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// Handle database errors specifically
		if err.Error() == "sql: no rows in result set" {
			ErrorResponse(w, r, apperrors.Unauthorized("Invalid credentials", err), h.logger)
			return
		}

		ErrorResponse(w, r, apperrors.Unauthorized("Invalid credentials", err), h.logger)
		return
	}

	RespondOK(w, map[string]string{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary Register new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body docs.RegisterRequest true "User registration information"
// @Success 201 {object} docs.Response{data=docs.MessageResponse}
// @Failure 400 {object} docs.ErrorResponse
// @Failure 409 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := DecodeJSONBody(w, r, &req); err != nil {
		ErrorResponse(w, r, err, h.logger)
		return
	}

	err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		// Check for specific error types
		if err.Error() == "email already exists" {
			ErrorResponse(w, r, apperrors.DuplicateEntity("A user with this email already exists", err), h.logger)
			return
		}

		// Handle database errors specifically
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			ErrorResponse(w, r, apperrors.DuplicateEntity("A user with this email already exists", err), h.logger)
			return
		}

		ErrorResponse(w, r, apperrors.InvalidInput("Registration failed", err), h.logger)
		return
	}

	RespondCreated(w, map[string]string{"message": "User registered successfully"})
}

// @Summary Refresh token
// @Description Refresh an expired access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body docs.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} docs.Response{data=docs.TokenResponse}
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := DecodeJSONBody(w, r, &req); err != nil {
		ErrorResponse(w, r, err, h.logger)
		return
	}

	if req.RefreshToken == "" {
		ErrorResponse(w, r, apperrors.InvalidInput("Refresh token is required", nil), h.logger)
		return
	}

	h.logger.Info().Msg("Refreshing token")

	// Call the service to refresh the token
	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		// Handle specific error types
		if strings.Contains(err.Error(), "token is revoked") ||
			strings.Contains(err.Error(), "invalid token") ||
			strings.Contains(err.Error(), "no rows in result set") {
			ErrorResponse(w, r, apperrors.Unauthorized("Invalid refresh token", err), h.logger)
			return
		}

		ErrorResponse(w, r, apperrors.Internal("Failed to refresh token", err), h.logger)
		return
	}

	RespondOK(w, map[string]string{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
		"message":       "Token refreshed successfully",
	})
}

// @Summary Logout user
// @Description Revoke a refresh token to log out a user
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body docs.RefreshTokenRequest true "Refresh token to revoke"
// @Success 200 {object} docs.Response{data=docs.MessageResponse}
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := DecodeJSONBody(w, r, &req); err != nil {
		ErrorResponse(w, r, err, h.logger)
		return
	}

	if req.RefreshToken == "" {
		ErrorResponse(w, r, apperrors.InvalidInput("Refresh token is required", nil), h.logger)
		return
	}

	h.logger.Info().Msg("Logging out user")

	// Call the service to revoke the refresh token
	err := h.authService.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		// Handle specific error types
		if strings.Contains(err.Error(), "token is revoked") ||
			strings.Contains(err.Error(), "invalid token") ||
			strings.Contains(err.Error(), "no rows in result set") {
			ErrorResponse(w, r, apperrors.Unauthorized("Invalid refresh token", err), h.logger)
			return
		}

		ErrorResponse(w, r, apperrors.Internal("Failed to logout", err), h.logger)
		return
	}

	RespondOK(w, map[string]string{"message": "Logged out successfully"})
}
