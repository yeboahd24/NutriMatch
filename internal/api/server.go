package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/api/handler"
	authMiddleware "github.com/yeboahd24/nutrimatch/internal/api/middleware/auth"
	errorsMiddleware "github.com/yeboahd24/nutrimatch/internal/api/middleware/errors"
	"github.com/yeboahd24/nutrimatch/internal/config"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres"
	"github.com/yeboahd24/nutrimatch/internal/repository/postgres/db"
	"github.com/yeboahd24/nutrimatch/internal/service"
	"github.com/yeboahd24/nutrimatch/pkg/auth"
)

// Server represents the API server
type Server struct {
	Router *chi.Mux
	Config *config.AppConfig
	Logger zerolog.Logger
	DB     *sql.DB
}

// NewServer creates a new API server
func NewServer(cfg *config.AppConfig, logger zerolog.Logger) (*Server, error) {
	// Connect to database
	db, err := postgres.NewDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Error handling middleware
	errorHandler := errorsMiddleware.NewErrorHandler(logger)
	r.Use(errorHandler.Middleware)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting
	r.Use(httprate.LimitByIP(
		cfg.Security.RateLimit,
		cfg.Security.RateLimitWindow,
	))

	// Create server
	server := &Server{
		Router: r,
		Config: cfg,
		Logger: logger,
		DB:     db,
	}

	// Register routes
	if err := server.registerRoutes(); err != nil {
		return nil, err
	}

	return server, nil
}

// registerRoutes registers all API routes
func (s *Server) registerRoutes() error {
	// Create repositories
	queries := db.New(s.DB)
	userRepo := postgres.NewUserRepository(queries)
	profileRepo := postgres.NewProfileRepository(queries)
	foodRepo := postgres.NewFoodRepository(queries)
	authRepo := postgres.NewAuthRepository(queries)
	referenceRepo := postgres.NewReferenceRepository(queries)

	// Create services
	passwordService := auth.NewPasswordService(s.Config.Security)
	jwtService := auth.NewJWTService(s.Config.JWT)
	userService := service.NewUserService(userRepo, authRepo, jwtService, passwordService, s.Logger)
	authService := service.NewAuthService(userRepo, authRepo, jwtService, passwordService, s.Logger)
	profileService := service.NewProfileService(profileRepo, userRepo, s.Logger)
	foodService := service.NewFoodService(foodRepo, s.Logger)
	recommendationService := service.NewRecommendationService(foodRepo, profileRepo, s.Logger)
	referenceService := service.NewReferenceService(referenceRepo, s.Logger)

	// Create handlers
	authHandler := handler.NewAuthHandler(authService, s.Logger)
	userHandler := handler.NewUserHandler(userService, s.Logger)
	profileHandler := handler.NewProfileHandler(profileService, s.Logger)
	foodHandler := handler.NewFoodHandler(foodService, s.Logger, s.Config.JWT)
	recommendationHandler := handler.NewRecommendationHandler(recommendationService, s.Logger)
	referenceHandler := handler.NewReferenceHandler(referenceService, s.Logger)

	// Public routes
	s.Router.Group(func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// API version
		r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("v1.0.0"))
		})

		// Auth routes
		r.Route("/api/v1/auth", authHandler.RegisterRoutes)

		// Food routes (public)
		r.Route("/api/v1/foods", foodHandler.RegisterRoutes)

		// Reference data routes (public)
		r.Route("/api/v1/reference", referenceHandler.RegisterRoutes)

		// Public user routes
		r.Route("/api/v1/users", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})
	})

	// Protected routes
	s.Router.Group(func(r chi.Router) {
		// Use auth middleware
		r.Use(authMiddleware.Middleware(s.Config.JWT))

		// Protected user routes
		r.Route("/api/v1/users/me", func(r chi.Router) {
			r.Get("/", userHandler.GetProfile)
			r.Put("/", userHandler.UpdateProfile)
		})

		// Profile routes
		r.Route("/api/v1/profiles", profileHandler.RegisterRoutes)

		// Recommendation routes
		r.Route("/api/v1/recommendations", recommendationHandler.RegisterRoutes)
	})

	return nil
}

// Close closes the server resources
func (s *Server) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
