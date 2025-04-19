package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AppError represents an application error
type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// StatusCode returns the HTTP status code for the error
func (e *AppError) StatusCode() int {
	switch e.Code {
	case ErrInvalidInput:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrDuplicateEntity:
		return http.StatusConflict
	case ErrDatabaseOperation, ErrInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WriteJSON writes the error as JSON to the response writer
func (e *AppError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode())

	response := struct {
		Error *AppError `json:"error"`
	}{
		Error: e,
	}

	json.NewEncoder(w).Encode(response)
}

// Common error codes
const (
	ErrInvalidInput      = "INVALID_INPUT"
	ErrUnauthorized      = "UNAUTHORIZED"
	ErrForbidden         = "FORBIDDEN"
	ErrNotFound          = "NOT_FOUND"
	ErrInternal          = "INTERNAL_ERROR"
	ErrDatabaseOperation = "DATABASE_ERROR"
	ErrDuplicateEntity   = "DUPLICATE_ENTITY"
)

// New creates a new AppError
func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewWithDetails creates a new AppError with details
func NewWithDetails(code, message string, details interface{}, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// InvalidInput creates a new invalid input error
func InvalidInput(message string, err error) *AppError {
	return New(ErrInvalidInput, message, err)
}

// InvalidInputWithDetails creates a new invalid input error with details
func InvalidInputWithDetails(message string, details interface{}, err error) *AppError {
	return NewWithDetails(ErrInvalidInput, message, details, err)
}

// Unauthorized creates a new unauthorized error
func Unauthorized(message string, err error) *AppError {
	return New(ErrUnauthorized, message, err)
}

// Forbidden creates a new forbidden error
func Forbidden(message string, err error) *AppError {
	return New(ErrForbidden, message, err)
}

// NotFound creates a new not found error
func NotFound(message string, err error) *AppError {
	return New(ErrNotFound, message, err)
}

// Internal creates a new internal error
func Internal(message string, err error) *AppError {
	return New(ErrInternal, message, err)
}

// DatabaseOperation creates a new database operation error
func DatabaseOperation(message string, err error) *AppError {
	return New(ErrDatabaseOperation, message, err)
}

// DuplicateEntity creates a new duplicate entity error
func DuplicateEntity(message string, err error) *AppError {
	return New(ErrDuplicateEntity, message, err)
}

// FromDatabaseError converts common database errors to AppError
func FromDatabaseError(err error, entityName string) *AppError {
	if err == nil {
		return nil
	}

	// Check for "no rows" error
	if err.Error() == "sql: no rows in result set" {
		return NotFound(fmt.Sprintf("%s not found", entityName), err)
	}

	// Check for unique constraint violation
	if strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "duplicate key") {
		return DuplicateEntity(fmt.Sprintf("%s already exists", entityName), err)
	}

	// Default to database operation error
	return DatabaseOperation(fmt.Sprintf("Database error occurred while accessing %s", entityName), err)
}
