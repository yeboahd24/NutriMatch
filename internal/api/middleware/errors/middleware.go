package errors

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	apperrors "github.com/yeboahd24/nutrimatch/pkg/errors"
)

// ErrorHandler is a middleware that handles errors
type ErrorHandler struct {
	logger zerolog.Logger
}

// NewErrorHandler creates a new error handler middleware
func NewErrorHandler(logger zerolog.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Middleware returns a middleware function that handles errors
func (h *ErrorHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer that can capture the response
		crw := &captureResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(crw, r)
	})
}

// HandleError handles an error and writes the appropriate response
func (h *ErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	// Check if the error is an AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// Log the error
		h.logger.Error().
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
	h.logger.Error().
		Err(err).
		Str("path", r.URL.Path).
		Str("method", r.Method).
		Msg("Unexpected error")

	// Write the error response
	internalErr.WriteJSON(w)
}

// captureResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (crw *captureResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Status returns the status code
func (crw *captureResponseWriter) Status() int {
	return crw.statusCode
}
