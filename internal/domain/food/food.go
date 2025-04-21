package food

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Food represents a food item in the system
type Food struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	AlternateNames     []string               `json:"alternate_names,omitempty"`
	Description        string                 `json:"description,omitempty"`
	FoodType           string                 `json:"food_type,omitempty"`
	Source             []map[string]string    `json:"source,omitempty"`
	Serving            map[string]interface{} `json:"serving,omitempty"`
	Nutrition100g      map[string]interface{} `json:"nutrition_100g,omitempty"`
	EAN13              string                 `json:"ean_13,omitempty"`
	Labels             []string               `json:"labels,omitempty"`
	PackageSize        map[string]interface{} `json:"package_size,omitempty"`
	Ingredients        string                 `json:"ingredients,omitempty"`
	IngredientAnalysis map[string]interface{} `json:"ingredient_analysis,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// FoodRating represents a user's rating of a food item
type FoodRating struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FoodID    string    `json:"food_id"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comments  string    `json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SavedFood represents a food item saved by a user
type SavedFood struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FoodID    string    `json:"food_id"`
	ListType  string    `json:"list_type" validate:"required,oneof=favorites shopping_list watch_list"`
	CreatedAt time.Time `json:"created_at"`
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

	// Rating methods
	CreateRating(rating *FoodRating) error
	UpdateRating(rating *FoodRating) error
	GetRating(userID uuid.UUID, foodID string) (*FoodRating, error)
	ListUserRatings(userID uuid.UUID, limit, offset int) ([]FoodRating, error)
	DeleteRating(userID uuid.UUID, foodID string) error

	// Saved food methods
	SaveFood(saved *SavedFood) error
	GetSavedFood(userID uuid.UUID, foodID string, listType string) (*SavedFood, error)
	ListSavedFoods(userID uuid.UUID, listType string, limit, offset int) ([]SavedFood, error)
	DeleteSavedFood(userID uuid.UUID, foodID string, listType string) error
}

// Service defines the interface for food business logic
type Service interface {
	// Existing methods
	Create(food *Food) error
	GetByID(id string) (*Food, error)
	GetByEAN13(ean13 string) (*Food, error)
	List(limit, offset int) ([]Food, error)
	ListByType(foodType string, limit, offset int) ([]Food, error)
	Search(query string, limit, offset int) ([]Food, error)
	Count() (int64, error)
	Delete(id string) error
	Import(filePath string) (int, error)

	// Rating methods
	RateFood(ctx context.Context, userID uuid.UUID, foodID string, rating int, comments string) (*FoodRating, error)
	UpdateRating(ctx context.Context, userID uuid.UUID, foodID string, rating int, comments string) (*FoodRating, error)
	GetUserRating(ctx context.Context, userID uuid.UUID, foodID string) (*FoodRating, error)
	ListUserRatings(ctx context.Context, userID uuid.UUID, limit, offset int) ([]FoodRating, error)
	DeleteRating(ctx context.Context, userID uuid.UUID, foodID string) error

	// Saved food methods
	SaveFood(ctx context.Context, userID uuid.UUID, foodID string, listType string) (*SavedFood, error)
	ListSavedFoods(ctx context.Context, userID uuid.UUID, listType string, limit, offset int) ([]SavedFood, error)
	RemoveSavedFood(ctx context.Context, userID uuid.UUID, foodID string, listType string) error
}
