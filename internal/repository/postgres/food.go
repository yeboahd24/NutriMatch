package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

type foodRepository struct {
	queries *db.Queries
}

func NewFoodRepository(queries *db.Queries) food.Repository {
	return &foodRepository{
		queries: queries,
	}
}

func (r *foodRepository) GetByID(id string) (*food.Food, error) {
	f, err := r.queries.GetFoodByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return mapDbFoodToDomain(&f), nil
}

func (r *foodRepository) GetByEAN13(ean13 string) (*food.Food, error) {
	f, err := r.queries.GetFoodByEAN13(context.Background(), sql.NullString{String: ean13, Valid: true})
	if err != nil {
		return nil, err
	}
	return mapDbFoodToDomain(&f), nil
}

func (r *foodRepository) List(limit, offset int) ([]food.Food, error) {
	foods, err := r.queries.ListFoods(context.Background(), db.ListFoodsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]food.Food, len(foods))
	for i, f := range foods {
		result[i] = *mapDbFoodToDomain(&f)
	}
	return result, nil
}

func (r *foodRepository) ListByType(foodType string, limit, offset int) ([]food.Food, error) {
	foods, err := r.queries.ListFoodsByType(context.Background(), db.ListFoodsByTypeParams{
		FoodType: sql.NullString{String: foodType, Valid: true},
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]food.Food, len(foods))
	for i, f := range foods {
		result[i] = *mapDbFoodToDomain(&f)
	}
	return result, nil
}

func (r *foodRepository) Search(query string, limit, offset int) ([]food.Food, error) {
	foods, err := r.queries.SearchFoodsByName(context.Background(), db.SearchFoodsByNameParams{
		Column1: sql.NullString{String: query, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]food.Food, len(foods))
	for i, f := range foods {
		result[i] = *mapDbFoodToDomain(&f)
	}
	return result, nil
}

func (r *foodRepository) Count() (int64, error) {
	return r.queries.CountFoods(context.Background())
}

func (r *foodRepository) Create(food *food.Food) error {
	alternateNames, _ := json.Marshal(food.AlternateNames)
	source, _ := json.Marshal(food.Source)
	serving, _ := json.Marshal(food.Serving)
	nutrition, _ := json.Marshal(food.Nutrition100g)
	labels, _ := json.Marshal(food.Labels)
	packageSize, _ := json.Marshal(food.PackageSize)
	ingredientAnalysis, _ := json.Marshal(food.IngredientAnalysis)

	_, err := r.queries.CreateFood(context.Background(), db.CreateFoodParams{
		ID:                 food.ID,
		Name:               food.Name,
		AlternateNames:     pqtype.NullRawMessage{RawMessage: alternateNames, Valid: true},
		Description:        sql.NullString{String: food.Description, Valid: food.Description != ""},
		FoodType:           sql.NullString{String: food.FoodType, Valid: food.FoodType != ""},
		Source:             pqtype.NullRawMessage{RawMessage: source, Valid: true},
		Serving:            pqtype.NullRawMessage{RawMessage: serving, Valid: true},
		Nutrition100g:      pqtype.NullRawMessage{RawMessage: nutrition, Valid: true},
		Ean13:              sql.NullString{String: food.EAN13, Valid: food.EAN13 != ""},
		Labels:             pqtype.NullRawMessage{RawMessage: labels, Valid: true},
		PackageSize:        pqtype.NullRawMessage{RawMessage: packageSize, Valid: true},
		Ingredients:        sql.NullString{String: food.Ingredients, Valid: food.Ingredients != ""},
		IngredientAnalysis: pqtype.NullRawMessage{RawMessage: ingredientAnalysis, Valid: true},
	})
	return err
}

func (r *foodRepository) Delete(id string) error {
	return r.queries.DeleteFood(context.Background(), id)
}

func (r *foodRepository) CreateRating(rating *food.FoodRating) error {
	result, err := r.queries.CreateFoodRating(context.Background(), db.CreateFoodRatingParams{
		UserID:   rating.UserID,
		FoodID:   rating.FoodID,
		Rating:   int16(rating.Rating),
		Comments: sql.NullString{String: rating.Comments, Valid: rating.Comments != ""},
	})
	if err != nil {
		return err
	}

	rating.ID = result.ID
	rating.CreatedAt = result.CreatedAt.Time
	rating.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *foodRepository) UpdateRating(rating *food.FoodRating) error {
	result, err := r.queries.UpdateFoodRating(context.Background(), db.UpdateFoodRatingParams{
		UserID:   rating.UserID,
		FoodID:   rating.FoodID,
		Rating:   int16(rating.Rating),
		Comments: sql.NullString{String: rating.Comments, Valid: rating.Comments != ""},
	})
	if err != nil {
		return err
	}

	rating.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *foodRepository) GetRating(userID uuid.UUID, foodID string) (*food.FoodRating, error) {
	result, err := r.queries.GetFoodRating(context.Background(), db.GetFoodRatingParams{
		UserID: userID,
		FoodID: foodID,
	})
	if err != nil {
		return nil, err
	}

	return mapDbRatingToDomain(&result), nil
}

func (r *foodRepository) ListUserRatings(userID uuid.UUID, limit, offset int) ([]food.FoodRating, error) {
	results, err := r.queries.ListUserRatings(context.Background(), db.ListUserRatingsParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	ratings := make([]food.FoodRating, len(results))
	for i, r := range results {
		ratings[i] = *mapDbRatingToDomain(&r)
	}
	return ratings, nil
}

func (r *foodRepository) DeleteRating(userID uuid.UUID, foodID string) error {
	return r.queries.DeleteFoodRating(context.Background(), db.DeleteFoodRatingParams{
		UserID: userID,
		FoodID: foodID,
	})
}

func (r *foodRepository) SaveFood(saved *food.SavedFood) error {
	result, err := r.queries.SaveFood(context.Background(), db.SaveFoodParams{
		UserID:   saved.UserID,
		FoodID:   saved.FoodID,
		ListType: saved.ListType,
	})
	if err != nil {
		return err
	}

	saved.ID = result.ID
	saved.CreatedAt = result.CreatedAt.Time
	return nil
}

func (r *foodRepository) GetSavedFood(userID uuid.UUID, foodID string, listType string) (*food.SavedFood, error) {
	result, err := r.queries.GetSavedFood(context.Background(), db.GetSavedFoodParams{
		UserID:   userID,
		FoodID:   foodID,
		ListType: listType,
	})
	if err != nil {
		return nil, err
	}

	return mapDbSavedFoodToDomain(&result), nil
}

func (r *foodRepository) ListSavedFoods(userID uuid.UUID, listType string, limit, offset int) ([]food.SavedFood, error) {
	results, err := r.queries.ListSavedFoods(context.Background(), db.ListSavedFoodsParams{
		UserID:   userID,
		ListType: listType,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	savedFoods := make([]food.SavedFood, len(results))
	for i, s := range results {
		savedFoods[i] = *mapDbSavedFoodToDomain(&s)
	}
	return savedFoods, nil
}

func (r *foodRepository) DeleteSavedFood(userID uuid.UUID, foodID string, listType string) error {
	return r.queries.DeleteSavedFood(context.Background(), db.DeleteSavedFoodParams{
		UserID:   userID,
		FoodID:   foodID,
		ListType: listType,
	})
}

func mapDbFoodToDomain(f *db.Food) *food.Food {
	var alternateNames []string
	var source []map[string]string
	var serving map[string]interface{}
	var nutrition map[string]interface{}
	var labels []string
	var packageSize map[string]interface{}
	var ingredientAnalysis map[string]interface{}

	json.Unmarshal(f.AlternateNames.RawMessage, &alternateNames)
	json.Unmarshal(f.Source.RawMessage, &source)
	json.Unmarshal(f.Serving.RawMessage, &serving)
	json.Unmarshal(f.Nutrition100g.RawMessage, &nutrition)
	json.Unmarshal(f.Labels.RawMessage, &labels)
	json.Unmarshal(f.PackageSize.RawMessage, &packageSize)
	json.Unmarshal(f.IngredientAnalysis.RawMessage, &ingredientAnalysis)

	return &food.Food{
		ID:                 f.ID,
		Name:               f.Name,
		AlternateNames:     alternateNames,
		Description:        f.Description.String,
		FoodType:           f.FoodType.String,
		Source:             source,
		Serving:            serving,
		Nutrition100g:      nutrition,
		EAN13:              f.Ean13.String,
		Labels:             labels,
		PackageSize:        packageSize,
		Ingredients:        f.Ingredients.String,
		IngredientAnalysis: ingredientAnalysis,
		CreatedAt:          f.CreatedAt.Time,
		UpdatedAt:          f.UpdatedAt.Time,
	}
}

func mapDbRatingToDomain(r *db.FoodRating) *food.FoodRating {
	return &food.FoodRating{
		ID:        r.ID,
		UserID:    r.UserID,
		FoodID:    r.FoodID,
		Rating:    int(r.Rating),
		Comments:  r.Comments.String,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func mapDbSavedFoodToDomain(s *db.UserSavedFood) *food.SavedFood {
	return &food.SavedFood{
		ID:        s.ID,
		UserID:    s.UserID,
		FoodID:    s.FoodID,
		ListType:  s.ListType,
		CreatedAt: s.CreatedAt.Time,
	}
}
