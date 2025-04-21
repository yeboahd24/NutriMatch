package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

type foodService struct {
	repo   food.Repository
	logger zerolog.Logger
}

func NewFoodService(
	repo food.Repository,
	logger zerolog.Logger,
) FoodService {
	return &foodService{
		repo:   repo,
		logger: logger,
	}
}

func (s *foodService) Create(food *food.Food) error {
	return s.repo.Create(food)
}

func (s *foodService) GetByID(id string) (*food.Food, error) {
	return s.repo.GetByID(id)
}

func (s *foodService) GetByEAN13(ean13 string) (*food.Food, error) {
	return s.repo.GetByEAN13(ean13)
}

func (s *foodService) List(limit, offset int) ([]food.Food, error) {
	return s.repo.List(limit, offset)
}

func (s *foodService) ListByType(foodType string, limit, offset int) ([]food.Food, error) {
	return s.repo.ListByType(foodType, limit, offset)
}

func (s *foodService) Search(query string, limit, offset int) ([]food.Food, error) {
	return s.repo.Search(query, limit, offset)
}

func (s *foodService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *foodService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *foodService) Import(filePath string) (int, error) {
	// Access underlying queries using type assertion
	if repo, ok := s.repo.(interface{ GetQueries() *db.Queries }); ok {
		importer := NewFoodImporter(repo.GetQueries(), s.logger)
		return importer.ImportFromTSV(filePath, 100) // Use batch size of 100
	}
	return 0, fmt.Errorf("repository does not support GetQueries")
}

// Adapter methods to implement the service.FoodService interface
func (s *foodService) GetFood(ctx context.Context, id string) (*food.Food, error) {
	return s.GetByID(id)
}

func (s *foodService) SearchFoods(ctx context.Context, query string, page, limit int) ([]food.Food, int, error) {
	offset := (page - 1) * limit
	foods, err := s.Search(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.Count()
	if err != nil {
		return foods, 0, err
	}

	return foods, int(count), nil
}

func (s *foodService) GetFoodsByCategory(ctx context.Context, category string, page, limit int) ([]food.Food, int, error) {
	offset := (page - 1) * limit
	foods, err := s.ListByType(category, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// This is a simplification - in a real implementation we would count only foods of this category
	count, err := s.Count()
	if err != nil {
		return foods, 0, err
	}

	return foods, int(count), nil
}

// Rating methods
func (s *foodService) RateFood(ctx context.Context, userID uuid.UUID, foodID string, rating int, comments string) (*food.FoodRating, error) {
	// Validate food exists
	if _, err := s.GetFood(ctx, foodID); err != nil {
		return nil, fmt.Errorf("food not found: %w", err)
	}

	// Create rating
	foodRating := &food.FoodRating{
		UserID:   userID,
		FoodID:   foodID,
		Rating:   rating,
		Comments: comments,
	}

	if err := s.repo.CreateRating(foodRating); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("user has already rated this food")
		}
		return nil, fmt.Errorf("failed to create rating: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("food_id", foodID).
		Int("rating", rating).
		Msg("Food rated successfully")

	return foodRating, nil
}

func (s *foodService) UpdateRating(ctx context.Context, userID uuid.UUID, foodID string, rating int, comments string) (*food.FoodRating, error) {
	// Get existing rating
	existing, err := s.repo.GetRating(userID, foodID)
	if err != nil {
		return nil, fmt.Errorf("rating not found: %w", err)
	}

	// Update rating
	existing.Rating = rating
	existing.Comments = comments

	if err := s.repo.UpdateRating(existing); err != nil {
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("food_id", foodID).
		Int("rating", rating).
		Msg("Rating updated successfully")

	return existing, nil
}

func (s *foodService) GetUserRating(ctx context.Context, userID uuid.UUID, foodID string) (*food.FoodRating, error) {
	rating, err := s.repo.GetRating(userID, foodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}
	return rating, nil
}

func (s *foodService) ListUserRatings(ctx context.Context, userID uuid.UUID, limit, offset int) ([]food.FoodRating, error) {
	if limit < 1 {
		limit = 10
	}
	ratings, err := s.repo.ListUserRatings(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list ratings: %w", err)
	}
	return ratings, nil
}

func (s *foodService) DeleteRating(ctx context.Context, userID uuid.UUID, foodID string) error {
	// Verify rating exists
	if _, err := s.repo.GetRating(userID, foodID); err != nil {
		return fmt.Errorf("rating not found: %w", err)
	}

	if err := s.repo.DeleteRating(userID, foodID); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("food_id", foodID).
		Msg("Rating deleted successfully")

	return nil
}

// Saved food methods
func (s *foodService) SaveFood(ctx context.Context, userID uuid.UUID, foodID string, listType string) (*food.SavedFood, error) {
	// Validate food exists
	if _, err := s.GetFood(ctx, foodID); err != nil {
		return nil, fmt.Errorf("food not found: %w", err)
	}

	savedFood := &food.SavedFood{
		UserID:   userID,
		FoodID:   foodID,
		ListType: listType,
	}

	if err := s.repo.SaveFood(savedFood); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("food already saved to this list")
		}
		return nil, fmt.Errorf("failed to save food: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("food_id", foodID).
		Str("list_type", listType).
		Msg("Food saved successfully")

	return savedFood, nil
}

func (s *foodService) ListSavedFoods(ctx context.Context, userID uuid.UUID, listType string, limit, offset int) ([]food.SavedFood, error) {
	if limit < 1 {
		limit = 10
	}
	saved, err := s.repo.ListSavedFoods(userID, listType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list saved foods: %w", err)
	}
	return saved, nil
}

func (s *foodService) RemoveSavedFood(ctx context.Context, userID uuid.UUID, foodID string, listType string) error {
	// Verify saved food exists
	if _, err := s.repo.GetSavedFood(userID, foodID, listType); err != nil {
		return fmt.Errorf("saved food not found: %w", err)
	}

	if err := s.repo.DeleteSavedFood(userID, foodID, listType); err != nil {
		return fmt.Errorf("failed to remove saved food: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("food_id", foodID).
		Str("list_type", listType).
		Msg("Saved food removed successfully")

	return nil
}
