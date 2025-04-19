package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/yeboahd24/nutrimatch/internal/config"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/migrate.go [up|down]")
	}
	
	command := os.Args[1]
	if command != "up" && command != "down" {
		log.Fatal("Command must be either 'up' or 'down'")
	}
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Connect to database
	db, err := postgres.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Get migrations path
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}
	
	// Run migrations
	if command == "up" {
		log.Println("Running migrations...")
		if err := postgres.RunMigrations(db, migrationsPath); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
	} else {
		log.Println("Rolling back last migration...")
		if err := postgres.RollbackMigration(db, migrationsPath); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Println("Rollback completed successfully")
	}
}
