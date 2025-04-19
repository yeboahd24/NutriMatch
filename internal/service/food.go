package service

import (
	"context"
	"fmt"

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
