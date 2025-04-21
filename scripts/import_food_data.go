package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sqlc-dev/pqtype"
	"github.com/yeboahd24/nutrimatch/internal/config"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
)

func importFoodsFromFile(foodsFilePath string) error {
	file, err := os.Open(foodsFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip header
	scanner.Scan()

	// Connect to database using the proper configuration
	cfg := config.DBConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "postgres",
		DBName:          "nutrimatch",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5,
	}

	pgdb, err := postgres.NewDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pgdb.Close()

	queries := db.New(pgdb)

	// Process each line
	lineCount := 0
	skippedCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) < 13 {
			fmt.Printf("Warning: Line %d has insufficient fields\n", lineCount)
			continue
		}

		// Parse fields
		id := fields[0]

		// Check if food already exists
		_, err := queries.GetFoodByID(context.Background(), id)
		if err == nil {
			// Food exists, skip it
			fmt.Printf("Skipping existing food with ID: %s\n", id)
			skippedCount++
			lineCount++
			continue
		} else if err != sql.ErrNoRows {
			// Unexpected error
			return fmt.Errorf("error checking food existence: %w", err)
		}

		name := fields[1]
		alternateNames := parseJSON(fields[2])
		description := fields[3]
		foodType := fields[4]
		source := parseJSON(fields[5])
		serving := parseJSON(fields[6])
		nutrition := parseNutrition(fields[7])
		ean13 := fields[8]
		labels := parseJSON(fields[9])
		packageSize := parseJSON(fields[10])
		ingredients := fields[11]
		ingredientAnalysis := parseJSON(fields[12])

		// Create food entry
		_, err = queries.CreateFood(context.Background(), db.CreateFoodParams{
			ID:                 id,
			Name:               name,
			AlternateNames:     pqtype.NullRawMessage{RawMessage: alternateNames, Valid: true},
			Description:        sql.NullString{String: description, Valid: description != ""},
			FoodType:           sql.NullString{String: foodType, Valid: foodType != ""},
			Source:             pqtype.NullRawMessage{RawMessage: source, Valid: true},
			Serving:            pqtype.NullRawMessage{RawMessage: serving, Valid: true},
			Nutrition100g:      pqtype.NullRawMessage{RawMessage: nutrition, Valid: true},
			Ean13:              sql.NullString{String: ean13, Valid: ean13 != ""},
			Labels:             pqtype.NullRawMessage{RawMessage: labels, Valid: true},
			PackageSize:        pqtype.NullRawMessage{RawMessage: packageSize, Valid: true},
			Ingredients:        sql.NullString{String: ingredients, Valid: ingredients != ""},
			IngredientAnalysis: pqtype.NullRawMessage{RawMessage: ingredientAnalysis, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to insert food %s at line %d: %w", id, lineCount, err)
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fmt.Printf("Successfully processed %d lines (skipped %d existing records)\n", lineCount, skippedCount)
	return nil
}

func parseJSON(s string) []byte {
	var data interface{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return []byte("{}")
	}
	bytes, _ := json.Marshal(data)
	return bytes
}

func parseNutrition(s string) []byte {
	var rawData map[string]string
	err := json.Unmarshal([]byte(s), &rawData)
	if err != nil {
		return []byte("{}")
	}

	// Convert string values to numbers
	nutrition := make(map[string]interface{})
	for key, value := range rawData {
		if value == "" {
			nutrition[key] = 0
			continue
		}

		// Try parsing as float first
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			nutrition[key] = floatVal
			continue
		}

		// If not a number, keep as string
		nutrition[key] = value
	}

	bytes, _ := json.Marshal(nutrition)
	return bytes
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run import.go <foods_file>")
		os.Exit(1)
	}

	err := importFoodsFromFile(os.Args[1])
	if err != nil {
		fmt.Printf("Import failed: %v\n", err)
		os.Exit(1)
	}
}
