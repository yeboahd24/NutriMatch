package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yeboahd24/nutrimatch/internal/config"
	"github.com/yeboahd24/nutrimatch/pkg/auth"
)

// contextKey is a custom type for context keys
type contextKey string

// Context keys
const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

// Middleware creates a new authentication middleware
func Middleware(cfg config.JWTConfig) func(http.Handler) http.Handler {
	jwtService := auth.NewJWTService(cfg)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}
			
			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}
			
			// Validate the token
			claims, err := jwtService.ValidateAccessToken(parts[1])
			if err != nil {
				if errors.Is(err, auth.ErrExpiredToken) {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				}
				if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, jwt.ErrSignatureInvalid) {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Failed to validate token", http.StatusInternalServerError)
				return
			}
			
			// Add user information to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			
			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID gets the user ID from the request context
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetEmail gets the email from the request context
func GetEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(EmailKey).(string)
	return email, ok
}
