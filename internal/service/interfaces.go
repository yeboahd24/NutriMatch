package service

import (
	"context"

	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
	"github.com/yeboahd24/nutrimatch/internal/domain/recommendation"
	"github.com/yeboahd24/nutrimatch/internal/domain/user"
)

// AuthService handles authentication operations
type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	Register(ctx context.Context, email, password, name string) error
	RefreshToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
}

// UserService handles user management operations
type UserService interface {
	GetProfile(ctx context.Context, userID string) (*user.User, error)
	UpdateProfile(ctx context.Context, userID, name, email string) error
}

// ProfileService handles user profile operations
type ProfileService interface {
	CreateProfile(ctx context.Context, userID string, age int, gender string, weight, height float64, goals, allergies, preferences []string, isDefault bool) (*profile.UserProfile, error)
	GetProfile(ctx context.Context, id string) (*profile.UserProfile, error)
	GetProfilesByUserID(ctx context.Context, userID string) ([]profile.UserProfile, error)
	UpdateProfile(ctx context.Context, id string, age int, gender string, weight, height float64, goals, allergies, preferences []string, isDefault bool) error
	DeleteProfile(ctx context.Context, id string) error
	GetAllProfiles(ctx context.Context) ([]*profile.UserProfile, error) // For debugging
}

// FoodService handles food-related operations
type FoodService interface {
	SearchFoods(ctx context.Context, query string, page, limit int) ([]food.Food, int, error)
	GetFood(ctx context.Context, id string) (*food.Food, error)
	GetFoodsByCategory(ctx context.Context, category string, page, limit int) ([]food.Food, int, error)
}

// RecommendationService handles food recommendation operations
type RecommendationService interface {
	GetDailyRecommendations(ctx context.Context, profileID string, limit int) ([]food.Food, error)
	GetMealPlanRecommendations(ctx context.Context, profileID string, days int) (*recommendation.MealPlan, error)
	GetFoodAlternatives(ctx context.Context, foodID string, limit int) ([]food.Food, error)
}
