package main

import (
	"context"
	"fmt"
	"identity-service/config"
	"identity-service/internal/domain"
	"identity-service/internal/handler"
	"identity-service/internal/middleware"
	"identity-service/internal/repository/postgres"
	"identity-service/internal/router"
	"identity-service/internal/service"
	"identity-service/pkg/database"
	"identity-service/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting Identity Service...")

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection
	db, err := database.GetDB(&cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.CloseDB()

	// Run database migrations
	if err := db.AutoMigrate(&domain.User{}, &domain.Address{}, &domain.Shop{}); err != nil {
		appLogger.Fatal("Failed to run migrations", zap.Error(err))
	}
	appLogger.Info("Database migrations completed")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	addressRepo := postgres.NewAddressRepository(db)
	shopRepo := postgres.NewShopRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, appLogger, cfg.JWT.Secret)
	userService := service.NewUserService(userRepo, appLogger)
	addressService := service.NewAddressService(addressRepo, appLogger)
	shopService := service.NewShopService(shopRepo, userRepo, appLogger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, appLogger)
	userHandler := handler.NewUserHandler(userService, appLogger)
	addressHandler := handler.NewAddressHandler(addressService, appLogger)
	shopHandler := handler.NewShopHandler(shopService, appLogger)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// Setup router
	router := router.SetupRouter(authHandler, userHandler, addressHandler, shopHandler, authMiddleware)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited gracefully")
}


