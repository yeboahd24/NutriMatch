package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api/middleware/auth"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
	"github.com/yeboahd24/nutrimatch/internal/service"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
	"github.com/yeboahd24/nutrimatch/pkg/response"
)

type UserHandler struct {
	BaseHandler
	userService service.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService service.UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(logger),
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get("/me", h.GetProfile)
	r.Put("/me", h.UpdateProfile)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Get("/user", h.GetUser)
	r.Put("/user", h.UpdateUser)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Use auth middleware helper to get user ID
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetProfile(r.Context(), userID.String())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user profile")
		response.Error(w, apperrors.Internal("Failed to get user profile", err))
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid user ID format", err))
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request payload", err))
		return
	}

	h.logger.Info().Str("user_id", userID.String()).Str("name", req.Name).Str("email", req.Email).Msg("User profile update requested")

	err = h.userService.UpdateProfile(r.Context(), userID.String(), req.Name, req.Email)
	if err != nil {
		if err == user.ErrEmailUpdateNotAllowed {
			response.Error(w, apperrors.InvalidInput("Email updates are not allowed for security reasons. Please contact support if you need to change your email.", err))
			return
		}
		if err.Error() == "sql: no rows in result set" {
			response.Error(w, apperrors.NotFound("user profile not found", err))
			return
		}
		response.Error(w, apperrors.Internal("Failed to update user profile", err))
		return
	}

	h.logger.Info().Str("user_id", userID.String()).Msg("User profile updated successfully")
	response.NoContent(w)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input user.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request payload", err))
		return
	}

	if err := h.validator.Struct(input); err != nil {
		response.Error(w, apperrors.InvalidInput("Validation failed", err))
		return
	}

	newUser, err := h.userService.CreateUser(r.Context(), &input)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			response.Error(w, apperrors.DuplicateEntity("A user with this email already exists", err))
			return
		}
		h.logger.Error().Err(err).Interface("input", input).Msg("Failed to register user")
		response.Error(w, apperrors.Internal("Failed to register user", err))
		return
	}

	response.JSON(w, http.StatusCreated, newUser)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request payload", err))
		return
	}

	if err := h.validator.Struct(input); err != nil {
		response.Error(w, apperrors.InvalidInput("Validation failed", err))
		return
	}

	accessToken, refreshToken, err := h.userService.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			response.Error(w, apperrors.Unauthorized("Invalid email or password", err))
			return
		}
		h.logger.Error().Err(err).Str("email", input.Email).Msg("Failed to login user")
		response.Error(w, apperrors.Internal("Failed to login user", err))
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid user ID format", err))
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID.String())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.Error(w, apperrors.NotFound("user not found", err))
			return
		}
		h.logger.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get user")
		response.Error(w, apperrors.Internal("Failed to get user", err))
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input user.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid request payload", err))
		return
	}

	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, apperrors.InvalidInput("Invalid user ID format", err))
		return
	}

	if err := h.validator.Struct(input); err != nil {
		response.Error(w, apperrors.InvalidInput("Validation failed", err))
		return
	}

	updatedUser, err := h.userService.UpdateUser(r.Context(), userID.String(), &input)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			response.Error(w, apperrors.DuplicateEntity("A user with this email already exists", err))
			return
		}
		h.logger.Error().Err(err).Interface("input", input).Msg("Failed to update user")
		response.Error(w, apperrors.Internal("Failed to update user", err))
		return
	}

	response.JSON(w, http.StatusOK, updatedUser)
}
