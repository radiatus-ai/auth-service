package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/radiatus-ai/auth-service/config"
	"github.com/radiatus-ai/auth-service/internal/auth"
	"github.com/radiatus-ai/auth-service/internal/middleware"
	"github.com/radiatus-ai/auth-service/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Run database migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)

	// Initialize services
	authService := auth.NewService(userRepo, orgRepo, cfg.JWTSecret, cfg.GoogleClientIDs, cfg.EmailWhitelist)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://dev.r7ai.net", "https://canvas.dev.r7ai.net"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	router.POST("/login/google", authHandler.LoginGoogle)
	router.POST("/api/verify-token", authHandler.VerifyToken)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		api.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"message": "You're authenticated!", "user_id": userID})
		})
	}

	// Start the server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
