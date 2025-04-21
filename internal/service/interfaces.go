package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/domain/profile"
	"github.com/yeboahd24/nutrimatch/internal/domain/recommendation"
	"github.com/yeboahd24/nutrimatch/internal/domain/reference"
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
	CreateUser(ctx context.Context, input *user.RegisterInput) (*user.User, error)
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	UpdateUser(ctx context.Context, userID string, input *user.UpdateUserInput) (*user.User, error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	GetByID(id uuid.UUID) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Update(user *user.User) error
	Register(email, password, firstName, lastName string) (*user.User, error)
	UpdatePassword(id uuid.UUID, currentPassword, newPassword string) error
	VerifyEmail(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

// ProfileService handles user profile operations
type ProfileService interface {
	CreateProfile(ctx context.Context, userID string, age int, gender string, weight, height float64, goals, allergies, preferences []string, isDefault bool) (*profile.UserProfile, error)
	GetProfile(ctx context.Context, id string) (*profile.UserProfile, error)
	GetProfilesByUserID(ctx context.Context, userID string) ([]profile.UserProfile, error)
	UpdateProfile(ctx context.Context, id string, age int, gender string, weight, height float64, goals, allergies, preferences []string, isDefault bool) error
	DeleteProfile(ctx context.Context, id string) error
	GetAllProfiles(ctx context.Context) ([]*profile.UserProfile, error)
}

// FoodService handles food-related operations
type FoodService interface {
	SearchFoods(ctx context.Context, query string, page, limit int) ([]food.Food, int, error)
	GetFood(ctx context.Context, id string) (*food.Food, error)
	GetFoodsByCategory(ctx context.Context, category string, page, limit int) ([]food.Food, int, error)
	Import(filePath string) (int, error)

	// Rating methods
	RateFood(ctx context.Context, userID uuid.UUID, foodID string, rating int, comments string) (*food.FoodRating, error)
	ListUserRatings(ctx context.Context, userID uuid.UUID, limit, offset int) ([]food.FoodRating, error)

	// Saved food methods
	SaveFood(ctx context.Context, userID uuid.UUID, foodID string, listType string) (*food.SavedFood, error)
	ListSavedFoods(ctx context.Context, userID uuid.UUID, listType string, limit, offset int) ([]food.SavedFood, error)
	RemoveSavedFood(ctx context.Context, userID uuid.UUID, foodID string, listType string) error
}

// RecommendationService handles food recommendation operations
type RecommendationService interface {
	GetRecommendations(userID uuid.UUID, req recommendation.RecommendationRequest) (*recommendation.RecommendationResponse, error)
	GetAlternatives(userID uuid.UUID, foodID string, limit int) ([]food.Food, error)
	GetDailyRecommendations(ctx context.Context, profileID string, limit int) ([]food.Food, error)
	GetMealPlanRecommendations(ctx context.Context, profileID string, days int) (*recommendation.MealPlan, error)
	GetFoodAlternatives(ctx context.Context, foodID string, limit int) ([]food.Food, error)
}

// ReferenceService handles reference data operations
type ReferenceService interface {
	GetAllergens(ctx context.Context) ([]reference.Allergen, error)
	GetHealthConditions(ctx context.Context) ([]reference.HealthCondition, error)
	GetDietaryPatterns(ctx context.Context) ([]reference.DietaryPattern, error)
}
