package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/import_foods.go <path-to-tsv-file>")
	}

	tsvFilePath := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(tsvFilePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", tsvFilePath)
	}

	// Connect to database
	db, err := sql.Open("pgx", "postgres://postgres:mesika@localhost:5432/nutrimatch?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Start import
	log.Printf("Starting import from %s", tsvFilePath)
	startTime := time.Now()

	count, err := importFoods(db, tsvFilePath, 100)
	if err != nil {
		log.Fatalf("Import failed: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Import completed successfully: %d foods imported in %v (%.2f items/sec)",
		count, duration, float64(count)/duration.Seconds())
}

func importFoods(db *sql.DB, filePath string, batchSize int) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // Set a larger buffer size

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

	// Prepare the insert statement
	stmt, err := db.Prepare(`
		INSERT INTO foods (
			id, name, alternate_names, description, food_type, source, 
			serving, nutrition_100g, ean_13, labels, package_size, 
			ingredients, ingredient_analysis, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			alternate_names = EXCLUDED.alternate_names,
			description = EXCLUDED.description,
			food_type = EXCLUDED.food_type,
			source = EXCLUDED.source,
			serving = EXCLUDED.serving,
			nutrition_100g = EXCLUDED.nutrition_100g,
			ean_13 = EXCLUDED.ean_13,
			labels = EXCLUDED.labels,
			package_size = EXCLUDED.package_size,
			ingredients = EXCLUDED.ingredients,
			ingredient_analysis = EXCLUDED.ingredient_analysis,
			updated_at = EXCLUDED.updated_at
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Start a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Import data
	count := 0
	batchCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		// Skip if we don't have enough fields
		if len(fields) < len(columns) {
			log.Printf("Skipping line with insufficient fields: %s", line)
			continue
		}

		// Get values from fields
		id := fields[columnMap["id"]]
		name := fields[columnMap["name"]]
		alternateNames := fields[columnMap["alternate_names"]]
		description := fields[columnMap["description"]]
		foodType := fields[columnMap["type"]]
		source := fields[columnMap["source"]]
		serving := fields[columnMap["serving"]]
		nutrition := fields[columnMap["nutrition_100g"]]
		ean13 := fields[columnMap["ean_13"]]
		labels := fields[columnMap["labels"]]
		packageSize := fields[columnMap["package_size"]]
		ingredients := fields[columnMap["ingredients"]]
		ingredientAnalysis := fields[columnMap["ingredient_analysis"]]

		// Validate JSON fields
		if !isValidJSON(alternateNames) {
			alternateNames = "[]"
		}
		if !isValidJSON(source) {
			source = "{}"
		}
		if !isValidJSON(serving) {
			serving = "{}"
		}
		if !isValidJSON(nutrition) {
			nutrition = "{}"
		}
		if !isValidJSON(labels) {
			labels = "[]"
		}
		if !isValidJSON(packageSize) {
			packageSize = "{}"
		}
		if !isValidJSON(ingredientAnalysis) {
			ingredientAnalysis = "{}"
		}

		now := time.Now()

		// Execute the insert
		_, err := tx.Stmt(stmt).Exec(
			id, name, alternateNames, description, foodType, source,
			serving, nutrition, ean13, labels, packageSize,
			ingredients, ingredientAnalysis, now, now,
		)
		if err != nil {
			log.Printf("Error inserting food %s: %v", id, err)
			continue
		}

		count++
		batchCount++

		// Commit and start a new transaction every batchSize records
		if batchCount >= batchSize {
			if err := tx.Commit(); err != nil {
				return count, fmt.Errorf("failed to commit transaction: %w", err)
			}

			// Start a new transaction
			tx, err = db.BeginTx(context.Background(), nil)
			if err != nil {
				return count, fmt.Errorf("failed to begin transaction: %w", err)
			}

			batchCount = 0
			log.Printf("Imported %d foods so far...", count)
		}
	}

	// Commit any remaining records
	if batchCount > 0 {
		if err := tx.Commit(); err != nil {
			return count, fmt.Errorf("failed to commit final transaction: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return count, fmt.Errorf("scanner error: %w", err)
	}

	return count, nil
}

func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
