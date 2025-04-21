package service

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/sqlc-dev/pqtype"
	"github.com/yeboahd24/nutrimatch/internal/domain/food"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

// FoodImporter handles importing food data from the OpenNutrition dataset
type FoodImporter struct {
	queries *db.Queries
	logger  zerolog.Logger
}

// NewFoodImporter creates a new food importer
func NewFoodImporter(queries *db.Queries, logger zerolog.Logger) *FoodImporter {
	return &FoodImporter{
		queries: queries,
		logger:  logger,
	}
}

// ImportFromTSV imports food data from a TSV file
func (i *FoodImporter) ImportFromTSV(filePath string, batchSize int) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read header line
	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to read header line")
	}
	header := scanner.Text()
	columns := strings.Split(header, "\t")

	// Create column index map
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}

	// Validate required columns
	requiredColumns := []string{
		"id", "name", "alternate_names", "description", "type",
		"source", "serving", "nutrition_100g", "ean_13",
		"labels", "package_size", "ingredients", "ingredient_analysis",
	}

	for _, col := range requiredColumns {
		if _, ok := columnMap[col]; !ok {
			return 0, fmt.Errorf("required column '%s' not found in TSV file", col)
		}
	}

	// Import data
	count := 0
	batch := make([]db.Food, 0, batchSize)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		// Skip if we don't have enough fields
		if len(fields) < len(columns) {
			i.logger.Warn().Str("line", line).Msg("Skipping line with insufficient fields")
			continue
		}

		// Parse food item
		food, err := i.parseFood(fields, columnMap)
		if err != nil {
			i.logger.Warn().Err(err).Str("line", line).Msg("Failed to parse food item")
			continue
		}

		// Marshal complex types to JSON
		alternateNamesJSON, err := json.Marshal(food.AlternateNames)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal alternate names")
			continue
		}

		sourceJSON, err := json.Marshal(food.Source)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal source")
			continue
		}

		servingJSON, err := json.Marshal(food.Serving)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal serving")
			continue
		}

		nutritionJSON, err := json.Marshal(food.Nutrition100g)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal nutrition")
			continue
		}

		labelsJSON, err := json.Marshal(food.Labels)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal labels")
			continue
		}

		packageSizeJSON, err := json.Marshal(food.PackageSize)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal package size")
			continue
		}

		ingredientAnalysisJSON, err := json.Marshal(food.IngredientAnalysis)
		if err != nil {
			i.logger.Warn().Err(err).Msg("Failed to marshal ingredient analysis")
			continue
		}

		batch = append(batch, db.Food{
			ID:                 food.ID,
			Name:               food.Name,
			AlternateNames:     pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(alternateNamesJSON)},
			Description:        sql.NullString{String: food.Description, Valid: food.Description != ""},
			FoodType:           sql.NullString{String: food.FoodType, Valid: food.FoodType != ""},
			Source:             pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(sourceJSON)},
			Serving:            pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(servingJSON)},
			Nutrition100g:      pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(nutritionJSON)},
			Ean13:              sql.NullString{String: food.EAN13, Valid: food.EAN13 != ""},
			Labels:             pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(labelsJSON)},
			PackageSize:        pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(packageSizeJSON)},
			Ingredients:        sql.NullString{String: food.Ingredients, Valid: food.Ingredients != ""},
			IngredientAnalysis: pqtype.NullRawMessage{Valid: true, RawMessage: json.RawMessage(ingredientAnalysisJSON)},
		})
		count++

		// Process batch
		if len(batch) >= batchSize {
			if err := i.processBatch(batch); err != nil {
				return count, fmt.Errorf("failed to process batch: %w", err)
			}
			batch = make([]db.Food, 0, batchSize)
		}
	}

	// Process remaining items
	if len(batch) > 0 {
		if err := i.processBatch(batch); err != nil {
			return count, fmt.Errorf("failed to process final batch: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return count, fmt.Errorf("scanner error: %w", err)
	}

	return count, nil
}

// parseFood parses a food item from TSV fields
func (i *FoodImporter) parseFood(fields []string, columnMap map[string]int) (food.Food, error) {
	var f food.Food
	var err error

	// Parse basic fields
	f.ID = fields[columnMap["id"]]
	f.Name = fields[columnMap["name"]]
	f.Description = fields[columnMap["description"]]
	f.FoodType = fields[columnMap["type"]]
	f.EAN13 = fields[columnMap["ean_13"]]
	f.Ingredients = fields[columnMap["ingredients"]]

	// Parse JSON fields
	if err = json.Unmarshal([]byte(fields[columnMap["alternate_names"]]), &f.AlternateNames); err != nil {
		return f, fmt.Errorf("failed to parse alternate_names: %w", err)
	}

	if err = json.Unmarshal([]byte(fields[columnMap["source"]]), &f.Source); err != nil {
		return f, fmt.Errorf("failed to parse source: %w", err)
	}

	if err = json.Unmarshal([]byte(fields[columnMap["serving"]]), &f.Serving); err != nil {
		return f, fmt.Errorf("failed to parse serving: %w", err)
	}

	if err = json.Unmarshal([]byte(fields[columnMap["nutrition_100g"]]), &f.Nutrition100g); err != nil {
		return f, fmt.Errorf("failed to parse nutrition_100g: %w", err)
	}

	if err = json.Unmarshal([]byte(fields[columnMap["labels"]]), &f.Labels); err != nil {
		return f, fmt.Errorf("failed to parse labels: %w", err)
	}

	if err = json.Unmarshal([]byte(fields[columnMap["package_size"]]), &f.PackageSize); err != nil {
		return f, fmt.Errorf("failed to parse package_size: %w", err)
	}

	if err = json.Unmarshal([]byte(fields[columnMap["ingredient_analysis"]]), &f.IngredientAnalysis); err != nil {
		return f, fmt.Errorf("failed to parse ingredient_analysis: %w", err)
	}

	return f, nil
}

// processBatch processes a batch of food items
func (i *FoodImporter) processBatch(foods []db.Food) error {
	ctx := context.Background()

	// Process all items in the transaction using the same queries instance
	for _, f := range foods {
		_, err := i.queries.CreateFood(ctx, db.CreateFoodParams{
			ID:                 f.ID,
			Name:               f.Name,
			AlternateNames:     f.AlternateNames,
			Description:        f.Description,
			FoodType:           f.FoodType,
			Source:             f.Source,
			Serving:            f.Serving,
			Nutrition100g:      f.Nutrition100g,
			Ean13:              f.Ean13,
			Labels:             f.Labels,
			PackageSize:        f.PackageSize,
			Ingredients:        f.Ingredients,
			IngredientAnalysis: f.IngredientAnalysis,
		})
		if err != nil {
			i.logger.Error().Err(err).Str("food_id", f.ID).Msg("Failed to insert food")
			return fmt.Errorf("failed to insert food %s: %w", f.ID, err)
		}
	}

	return nil
}
