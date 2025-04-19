package main

import (
	"log"
	"os"
	"time"

	"github.com/yeboahd24/nutrimatch/internal/config"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres"
	"github.com/yeboahd24/nutrimatch/internal/service"
	"github.com/yeboahd24/nutrimatch/pkg/logger"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/import.go <path-to-tsv-file>")
	}

	tsvFilePath := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(tsvFilePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", tsvFilePath)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	l := logger.New(cfg.Logging)

	// Connect to database
	db, err := postgres.NewDB(cfg.Database)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize sqlc queries
	queries := db.New(db)

	// Create importer
	importer := service.NewFoodImporter(queries, l)

	// Start import
	l.Info().Str("file", tsvFilePath).Msg("Starting import")
	startTime := time.Now()

	count, err := importer.ImportFromTSV(tsvFilePath, 100)
	if err != nil {
		l.Fatal().Err(err).Msg("Import failed")
	}

	duration := time.Since(startTime)
	l.Info().
		Int("count", count).
		Dur("duration", duration).
		Float64("items_per_second", float64(count)/duration.Seconds()).
		Msg("Import completed successfully")
}
