package reference

import (
	"context"
	"time"
)

// Allergen represents a food allergen
type Allergen struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	CommonNames map[string]interface{} `json:"common_names,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// HealthCondition represents a health condition with dietary implications
type HealthCondition struct {
	ID                      int                    `json:"id"`
	Name                    string                 `json:"name"`
	Description             string                 `json:"description,omitempty"`
	NutrientRestrictions    map[string]interface{} `json:"nutrient_restrictions,omitempty"`
	NutrientRecommendations map[string]interface{} `json:"nutrient_recommendations,omitempty"`
	CreatedAt               time.Time              `json:"created_at"`
}

// DietaryPattern represents a dietary pattern (e.g., vegetarian, vegan, etc.)
type DietaryPattern struct {
	Name            string   `json:"name"`
	Description     string   `json:"description,omitempty"`
	Restrictions    []string `json:"restrictions"`
	Recommendations []string `json:"recommendations"`
}

// Repository defines the interface for reference data access
type Repository interface {
	GetAllergens(ctx context.Context) ([]Allergen, error)
	GetHealthConditions(ctx context.Context) ([]HealthCondition, error)
	GetDietaryPatterns(ctx context.Context) ([]DietaryPattern, error)
}
