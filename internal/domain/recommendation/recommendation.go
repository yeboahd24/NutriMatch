package recommendation

import (
	"github.com/google/uuid"
	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
)

// Rule represents a filtering rule for food recommendations
type Rule struct {
	Type      string      `json:"type"`      // e.g., "allergen", "nutrient", "preference"
	Operation string      `json:"operation"` // e.g., "exclude", "include", "max", "min"
	Target    string      `json:"target"`    // e.g., "peanuts", "sodium", "calories"
	Value     interface{} `json:"value"`     // The threshold value for the rule
	Priority  int         `json:"priority"`  // Rule priority (higher = more important)
}

// Meal represents a single meal with recommended foods
type Meal struct {
	Type  string      `json:"type"` // breakfast, lunch, dinner, snack
	Foods []food.Food `json:"foods"`
}

// DailyPlan represents a full day of meal recommendations
type DailyPlan struct {
	Date  string `json:"date"`
	Meals []Meal `json:"meals"`
}

// MealPlan represents a complete meal plan for multiple days
type MealPlan struct {
	ProfileID string      `json:"profile_id"`
	Days      []DailyPlan `json:"days"`
	TotalDays int         `json:"total_days"`
}

// RecommendationRequest represents a request for food recommendations
type RecommendationRequest struct {
	ProfileID   *uuid.UUID `json:"profile_id,omitempty"`   // Optional: use specific profile
	CustomRules []Rule     `json:"custom_rules,omitempty"` // Optional: additional rules
	Limit       int        `json:"limit,omitempty"`        // Optional: limit results
	Offset      int        `json:"offset,omitempty"`       // Optional: pagination offset
}

// RecommendationResponse represents a response with food recommendations
type RecommendationResponse struct {
	Foods        []food.Food `json:"foods"`
	TotalCount   int         `json:"total_count"`
	AppliedRules []Rule      `json:"applied_rules"`
}

// Service defines the interface for recommendation business logic
type Service interface {
	GetRecommendations(userID uuid.UUID, req RecommendationRequest) (*RecommendationResponse, error)
	GetAlternatives(userID uuid.UUID, foodID string, limit int) ([]food.Food, error)
	GenerateRulesFromProfile(profile *profile.UserProfile) ([]Rule, error)
}
