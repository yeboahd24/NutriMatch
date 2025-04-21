package docs

import (
	"time"

	"github.com/google/uuid"
)

// This file contains model definitions for Swagger documentation.
// These models are used only for documentation purposes and don't affect the actual code.

// Auth Models

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"John Doe"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// TokenResponse represents the response for successful authentication
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// Profile Models

// ProfileRequest represents the request body for creating or updating a profile
type ProfileRequest struct {
	Age         int      `json:"age" example:"30"`
	Gender      string   `json:"gender" example:"male"`
	Weight      float64  `json:"weight" example:"75.5"`
	Height      float64  `json:"height" example:"180"`
	Goals       []string `json:"goals" example:"[\"weight_loss\", \"muscle_gain\"]"`
	Allergies   []string `json:"allergies" example:"[\"peanuts\", \"shellfish\"]"`
	Preferences []string `json:"preferences" example:"[\"vegetarian\", \"low_carb\"]"`
	IsDefault   bool     `json:"is_default" example:"true"`
}

// ProfileResponse represents a user profile in the API response
type ProfileResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Weight      float64   `json:"weight"`
	Height      float64   `json:"height"`
	Goals       []string  `json:"goals"`
	Allergies   []string  `json:"allergies"`
	Preferences []string  `json:"preferences"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Recommendation Models

// RecommendationRequest represents the request body for filtering recommendations
type RecommendationRequest struct {
	ProfileID   string   `json:"profile_id,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	MaxCalories int      `json:"max_calories,omitempty"`
	MinProtein  float64  `json:"min_protein,omitempty"`
	MaxFat      float64  `json:"max_fat,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	Offset      int      `json:"offset,omitempty"`
}

// RecommendationResponse represents the response for food recommendations
type RecommendationResponse struct {
	Foods        []FoodResponse `json:"recommendations"`
	TotalCount   int            `json:"total_count"`
	AppliedRules []string       `json:"applied_rules"`
	Pagination   struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"pagination"`
}

// AlternativesResponse represents the response for food alternatives
type AlternativesResponse struct {
	Alternatives []FoodResponse `json:"alternatives"`
	FoodID       string         `json:"food_id"`
	Limit        int            `json:"limit"`
}

// Reference Models

// ReferenceItem represents a reference data item
type ReferenceItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// FoodResponse represents a food item in the API response
type FoodResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Calories      float64 `json:"calories"`
	Protein       float64 `json:"protein"`
	Carbohydrates float64 `json:"carbohydrates"`
	Fat           float64 `json:"fat"`
	Fiber         float64 `json:"fiber"`
	Sugar         float64 `json:"sugar"`
	Sodium        float64 `json:"sodium"`
	ImageURL      string  `json:"image_url,omitempty"`
}

// FoodDetailResponse represents detailed information about a food item
type FoodDetailResponse struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Category            string                 `json:"category"`
	Calories            float64                `json:"calories"`
	Protein             float64                `json:"protein"`
	Carbohydrates       float64                `json:"carbohydrates"`
	Fat                 float64                `json:"fat"`
	Fiber               float64                `json:"fiber"`
	Sugar               float64                `json:"sugar"`
	Sodium              float64                `json:"sodium"`
	Ingredients         string                 `json:"ingredients,omitempty"`
	AllergenInfo        string                 `json:"allergen_info,omitempty"`
	ServingSize         string                 `json:"serving_size"`
	ServingSizeUnit     string                 `json:"serving_size_unit"`
	NutritionPerServing map[string]float64     `json:"nutrition_per_serving"`
	ImageURL            string                 `json:"image_url,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// RatingResponse represents a user's rating for a food item
type RatingResponse struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FoodID    string    `json:"food_id"`
	Rating    int       `json:"rating"`
	Comments  string    `json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// RatingRequest represents the request body for rating a food
type RatingRequest struct {
	Rating   int    `json:"rating" validate:"required,min=1,max=5"`
	Comments string `json:"comments,omitempty"`
}
