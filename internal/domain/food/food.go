package food

import (
	"time"
)

// Food represents a food item in the system
type Food struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	AlternateNames    []string               `json:"alternate_names,omitempty"`
	Description       string                 `json:"description,omitempty"`
	FoodType          string                 `json:"food_type,omitempty"`
	Source            []map[string]string    `json:"source,omitempty"`
	Serving           map[string]interface{} `json:"serving,omitempty"`
	Nutrition100g     map[string]interface{} `json:"nutrition_100g,omitempty"`
	EAN13             string                 `json:"ean_13,omitempty"`
	Labels            []string               `json:"labels,omitempty"`
	PackageSize       map[string]interface{} `json:"package_size,omitempty"`
	Ingredients       string                 `json:"ingredients,omitempty"`
	IngredientAnalysis map[string]interface{} `json:"ingredient_analysis,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// Repository defines the interface for food data access
type Repository interface {
	Create(food *Food) error
	GetByID(id string) (*Food, error)
	GetByEAN13(ean13 string) (*Food, error)
	List(limit, offset int) ([]Food, error)
	ListByType(foodType string, limit, offset int) ([]Food, error)
	Search(query string, limit, offset int) ([]Food, error)
	Count() (int64, error)
	Delete(id string) error
}

// Service defines the interface for food business logic
type Service interface {
	Create(food *Food) error
	GetByID(id string) (*Food, error)
	GetByEAN13(ean13 string) (*Food, error)
	List(limit, offset int) ([]Food, error)
	ListByType(foodType string, limit, offset int) ([]Food, error)
	Search(query string, limit, offset int) ([]Food, error)
	Count() (int64, error)
	Delete(id string) error
	Import(filePath string) (int, error)
}
