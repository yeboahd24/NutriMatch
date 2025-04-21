package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response represents a standardized API response
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status: http.StatusText(status),
		Data:   data,
	})
}

// JSONWithMeta sends a JSON response with metadata
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status: http.StatusText(status),
		Data:   data,
		Meta:   meta,
	})
}

// Error sends an error response
func Error(w http.ResponseWriter, err error) {
	var status int
	var message string

	switch e := err.(type) {
	case *ErrorResponse:
		status = e.Status
		message = e.Message
	default:
		status = http.StatusInternalServerError
		message = "Internal Server Error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status: http.StatusText(status),
		Error:  message,
	})
}

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error implements the error interface for ErrorResponse
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s (status: %d)", e.Message, e.Status)
}

// PaginationMeta creates pagination metadata
func PaginationMeta(page, limit, total int) map[string]interface{} {
	return map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	}
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
