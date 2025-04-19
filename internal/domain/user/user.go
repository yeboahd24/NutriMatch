package user

import (
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
}
