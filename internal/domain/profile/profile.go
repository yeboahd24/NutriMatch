package profile

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrUnauthorized = errors.New("unauthorized access to profile")
)

// UserProfile represents a user's nutritional profile
type UserProfile struct {
	ID                      uuid.UUID `json:"id"`
	UserID                  uuid.UUID `json:"user_id"`
	ProfileName             string    `json:"profile_name"`
	IsDefault               bool      `json:"is_default"`
	HealthConditions        []string  `json:"health_conditions"`
	DietaryRestrictions     []string  `json:"dietary_restrictions"`
	Allergens               []string  `json:"allergens"`
	GoalType                string    `json:"goal_type,omitempty"`
	CalorieTarget           int       `json:"calorie_target,omitempty"`
	MacronutrientPreference string    `json:"macronutrient_preference,omitempty"`
	DislikedFoods           []string  `json:"disliked_foods"`
	PreferredFoods          []string  `json:"preferred_foods"`
	CuisinePreferences      []string  `json:"cuisine_preferences"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// Repository defines the interface for profile data access
type Repository interface {
	Create(profile *UserProfile) error
	GetByID(id uuid.UUID) (*UserProfile, error)
	GetByUserID(userID uuid.UUID) ([]UserProfile, error)
	GetDefaultByUserID(userID uuid.UUID) (*UserProfile, error)
	Update(profile *UserProfile) error
	SetAsDefault(id uuid.UUID, userID uuid.UUID) error
	Delete(id uuid.UUID, userID uuid.UUID) error
	GetAll() ([]*UserProfile, error) // For debugging
}

// Service defines the interface for profile business logic
type Service interface {
	Create(profile *UserProfile) error
	GetByID(id uuid.UUID, userID uuid.UUID) (*UserProfile, error)
	GetByUserID(userID uuid.UUID) ([]UserProfile, error)
	GetDefaultByUserID(userID uuid.UUID) (*UserProfile, error)
	Update(profile *UserProfile) error
	SetAsDefault(id uuid.UUID, userID uuid.UUID) error
	Delete(id uuid.UUID, userID uuid.UUID) error
}
