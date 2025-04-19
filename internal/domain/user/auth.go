package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidPassword    = errors.New("invalid password")
)

// Credentials represents user login credentials
type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        uuid.UUID  `json:"-"`
	UserID    uuid.UUID  `json:"-"`
	Token     string     `json:"-"`
	ExpiresAt time.Time  `json:"-"`
	CreatedAt time.Time  `json:"-"`
	Revoked   bool       `json:"-"`
	RevokedAt *time.Time `json:"-"`
}

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(req RegisterRequest) (*User, error)
	Login(creds Credentials) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	Logout(refreshToken string) error
	LogoutAll(userID uuid.UUID) error
	VerifyPassword(hashedPassword, password string) bool
	HashPassword(password string) (string, error)
}

// AuthRepository defines the interface for auth data access
type AuthRepository interface {
	StoreRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error
	GetRefreshToken(token string) (*RefreshToken, error)
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID uuid.UUID) error
	DeleteExpiredRefreshTokens() error
}
