package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrEmailUpdateNotAllowed = errors.New("email updates are not allowed for security reasons")
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	FirstName     string     `json:"first_name,omitempty"`
	LastName      string     `json:"last_name,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Gender        string     `json:"gender,omitempty"`
	HeightCm      float64    `json:"height_cm,omitempty"`
	WeightKg      float64    `json:"weight_kg,omitempty"`
	ActivityLevel string     `json:"activity_level,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	AccountStatus string     `json:"account_status"`
	EmailVerified bool       `json:"email_verified"`
	MFAEnabled    bool       `json:"mfa_enabled"`
	MFASecret     string     `json:"-"`
}

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// UpdateUserInput represents the input for updating user information
type UpdateUserInput struct {
	FirstName     *string    `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName      *string    `json:"last_name,omitempty" validate:"omitempty,min=1"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Gender        *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	HeightCm      *float64   `json:"height_cm,omitempty" validate:"omitempty,gt=0"`
	WeightKg      *float64   `json:"weight_kg,omitempty" validate:"omitempty,gt=0"`
	ActivityLevel *string    `json:"activity_level,omitempty" validate:"omitempty,oneof=sedentary light moderate very_active"`
}

// Repository defines the interface for user data access
type Repository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	UpdatePassword(id uuid.UUID, passwordHash string) error
	UpdateLastLogin(id uuid.UUID) error
	UpdateEmailVerification(id uuid.UUID, verified bool) error
	Delete(id uuid.UUID) error
}

// Service defines the interface for user business logic
type Service interface {
	Register(email, password, firstName, lastName string) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	UpdatePassword(id uuid.UUID, currentPassword, newPassword string) error
	VerifyEmail(id uuid.UUID) error
	Delete(id uuid.UUID) error
	CreateUser(ctx context.Context, input *RegisterInput) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, userID string, input *UpdateUserInput) (*User, error)
	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, userID, name, email string) error
	Authenticate(ctx context.Context, email, password string) (string, error)
}
