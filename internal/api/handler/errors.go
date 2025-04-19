package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
)

// ErrorResponse handles errors in handlers
func ErrorResponse(w http.ResponseWriter, r *http.Request, err error, logger zerolog.Logger) {
	// Check if the error is an AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// Log the error
		logger.Error().
			Str("code", appErr.Code).
			Str("message", appErr.Message).
			Err(appErr.Err).
			Str("path", r.URL.Path).
			Str("method", r.Method).
			Msg("Application error")

		// Write the error response
		appErr.WriteJSON(w)
		return
	}

	// If it's not an AppError, create a generic internal error
	internalErr := apperrors.Internal("An unexpected error occurred", err)

	// Log the error
	logger.Error().
		Err(err).
		Str("path", r.URL.Path).
		Str("method", r.Method).
		Msg("Unexpected error")

	// Write the error response
	internalErr.WriteJSON(w)
}

// DecodeJSONBody decodes a JSON request body into the provided struct
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperrors.InvalidInput("Invalid request body", err)
	}
	return nil
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RespondCreated sends a 201 Created response with the provided data
func RespondCreated(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusCreated, data)
}

// RespondOK sends a 200 OK response with the provided data
func RespondOK(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, data)
}

// RespondNoContent sends a 204 No Content response
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// HandleDatabaseError handles common database errors and returns an appropriate AppError
func HandleDatabaseError(w http.ResponseWriter, r *http.Request, err error, entityName string, logger zerolog.Logger) {
	appErr := apperrors.FromDatabaseError(err, entityName)
	ErrorResponse(w, r, appErr, logger)
}
