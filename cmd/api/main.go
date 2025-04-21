// Package main is the entry point for the NutriMatch API server.
//
// @title NutriMatch API
// @version 1.0
// @description A personalized nutrition recommendation system built on the OpenNutrition dataset
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.nutrimatch.com/support
// @contact.email support@nutrimatch.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api"
	"github.com/yeboahd24/nutrimatch/internal/config"

	_ "github.com/yeboahd24/nutrimatch/docs"              // Import generated Swagger docs
	_ "github.com/yeboahd24/nutrimatch/internal/api/docs" // Import Swagger model definitions
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if cfg.Logging.Format == "pretty" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// Create and initialize server
	server, err := api.NewServer(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create server")
	}
	defer server.Close()

	// Create server with timeouts
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      server.Router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info().Int("port", cfg.Server.Port).Msg("Starting server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown gracefully
	logger.Info().Msg("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited properly")
}
