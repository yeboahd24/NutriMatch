package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/migrate_simple.go [up|down]")
	}
	
	command := os.Args[1]
	if command != "up" && command != "down" {
		log.Fatal("Command must be either 'up' or 'down'")
	}
	
	// Database connection string
	dbHost := getEnv("NUTRIMATCH_DATABASE_HOST", "localhost")
	dbPort := getEnv("NUTRIMATCH_DATABASE_PORT", "5432")
	dbUser := getEnv("NUTRIMATCH_DATABASE_USER", "postgres")
	dbPass := getEnv("NUTRIMATCH_DATABASE_PASSWORD", "postgres")
	dbName := getEnv("NUTRIMATCH_DATABASE_DBNAME", "postgres")
	dbSSLMode := getEnv("NUTRIMATCH_DATABASE_SSLMODE", "disable")
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)
	
	// Connect to database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	
	// Get migrations path
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}
	
	// Create migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}
	
	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	
	// Run migrations
	if command == "up" {
		log.Println("Running migrations...")
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to apply")
			} else {
				log.Fatalf("Failed to run migrations: %v", err)
			}
		} else {
			log.Println("Migrations completed successfully")
		}
	} else {
		log.Println("Rolling back last migration...")
		if err := m.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to roll back")
			} else {
				log.Fatalf("Failed to rollback migration: %v", err)
			}
		} else {
			log.Println("Rollback completed successfully")
		}
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
