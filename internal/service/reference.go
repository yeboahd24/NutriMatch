package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/domain/reference"
)

type referenceService struct {
	repo   reference.Repository
	logger zerolog.Logger
}

func NewReferenceService(
	repo reference.Repository,
	logger zerolog.Logger,
) ReferenceService {
	return &referenceService{
		repo:   repo,
		logger: logger,
	}
}

func (s *referenceService) GetAllergens(ctx context.Context) ([]reference.Allergen, error) {
	s.logger.Debug().Msg("Getting allergens from repository")
	return s.repo.GetAllergens(ctx)
}

func (s *referenceService) GetHealthConditions(ctx context.Context) ([]reference.HealthCondition, error) {
	s.logger.Debug().Msg("Getting health conditions from repository")
	return s.repo.GetHealthConditions(ctx)
}

func (s *referenceService) GetDietaryPatterns(ctx context.Context) ([]reference.DietaryPattern, error) {
	s.logger.Debug().Msg("Getting dietary patterns from repository")
	return s.repo.GetDietaryPatterns(ctx)
}
