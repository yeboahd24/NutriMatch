package postgres

import (
	"context"
	"encoding/json"

	"github.com/yeboahd24/nutrimatch/internal/domain/reference"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

type referenceRepository struct {
	queries *db.Queries
}

func NewReferenceRepository(queries *db.Queries) reference.Repository {
	return &referenceRepository{
		queries: queries,
	}
}

func (r *referenceRepository) GetAllergens(ctx context.Context) ([]reference.Allergen, error) {
	allergens, err := r.queries.ListAllergens(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]reference.Allergen, len(allergens))
	for i, a := range allergens {
		var commonNames map[string]interface{}
		if err := json.Unmarshal(a.CommonNames.RawMessage, &commonNames); err != nil {
			return nil, err
		}

		result[i] = reference.Allergen{
			ID:          int(a.ID),
			Name:        a.Name,
			Description: a.Description.String,
			CommonNames: commonNames,
			CreatedAt:   a.CreatedAt.Time,
		}
	}

	return result, nil
}

func (r *referenceRepository) GetHealthConditions(ctx context.Context) ([]reference.HealthCondition, error) {
	conditions, err := r.queries.ListHealthConditions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]reference.HealthCondition, len(conditions))
	for i, c := range conditions {
		var restrictions, recommendations map[string]interface{}
		if err := json.Unmarshal(c.NutrientRestrictions.RawMessage, &restrictions); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(c.NutrientRecommendations.RawMessage, &recommendations); err != nil {
			return nil, err
		}

		result[i] = reference.HealthCondition{
			ID:                      int(c.ID),
			Name:                    c.Name,
			Description:             c.Description.String,
			NutrientRestrictions:    restrictions,
			NutrientRecommendations: recommendations,
			CreatedAt:               c.CreatedAt.Time,
		}
	}

	return result, nil
}

func (r *referenceRepository) GetDietaryPatterns(ctx context.Context) ([]reference.DietaryPattern, error) {
	// For now, return a static list of dietary patterns
	// In the future, this could be moved to the database
	return []reference.DietaryPattern{
		{
			Name:        "vegetarian",
			Description: "A diet that excludes meat and fish but includes other animal products",
			Restrictions: []string{
				"meat",
				"fish",
				"seafood",
			},
			Recommendations: []string{
				"legumes",
				"nuts",
				"seeds",
				"dairy",
				"eggs",
			},
		},
		{
			Name:        "vegan",
			Description: "A diet that excludes all animal products",
			Restrictions: []string{
				"meat",
				"fish",
				"seafood",
				"dairy",
				"eggs",
				"honey",
			},
			Recommendations: []string{
				"legumes",
				"nuts",
				"seeds",
				"whole_grains",
				"fruits",
				"vegetables",
			},
		},
		{
			Name:        "pescatarian",
			Description: "A diet that includes fish but excludes other meats",
			Restrictions: []string{
				"meat",
				"poultry",
			},
			Recommendations: []string{
				"fish",
				"seafood",
				"vegetables",
				"fruits",
				"whole_grains",
			},
		},
	}, nil
}
