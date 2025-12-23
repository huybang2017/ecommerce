package router

import (
	"api-gateway/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter configures all API Gateway routes
func SetupRouter(
	gatewayHandler *handler.GatewayHandler,
	cfg *config.Config,
	logger *zap.Logger,
) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	if len(cfg.CORS.AllowedOrigins) > 0 {
		corsConfig := cors.Config{
			AllowOrigins:     cfg.CORS.AllowedOrigins,
			AllowMethods:     cfg.CORS.AllowedMethods,
			AllowHeaders:     cfg.CORS.AllowedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           cfg.CORS.MaxAge,
		}
		router.Use(cors.New(corsConfig))
	} else {
		// Default CORS config if not specified
		router.Use(cors.Default())
	}

	// Request logging middleware
	router.Use(middleware.RequestLoggingMiddleware(logger))
	router.Use(middleware.ErrorLoggingMiddleware(logger))

	// Rate limiting middleware
	router.Use(middleware.RateLimitMiddleware(&cfg.RateLimit, logger))

	// Health check endpoint (no auth required)
	router.GET("/health", gatewayHandler.HealthCheck)
	router.GET("/api/gateway/health", gatewayHandler.HealthCheck)

		// API routes - all requests go through the gateway
		api := router.Group("/api")
		{
			v1 := api.Group("/v1")
			{
				// Product service routes
				products := v1.Group("/products")
				{
					// Public routes (no auth required)
					products.GET("", gatewayHandler.ProxyRequest)
					products.GET("/:id", gatewayHandler.ProxyRequest)
					products.GET("/search", gatewayHandler.ProxyRequest)
					products.POST("", gatewayHandler.ProxyRequest)

					// Protected routes (auth required)
					protected := products.Group("")
					protected.Use(middleware.AuthMiddleware(&cfg.JWT, logger))
					{
						protected.PUT("/:id", gatewayHandler.ProxyRequest)
						protected.PATCH("/:id", gatewayHandler.ProxyRequest)
						protected.PATCH("/:id/inventory", gatewayHandler.ProxyRequest)
						protected.DELETE("/:id", gatewayHandler.ProxyRequest)
					}
				}

				// Identity service routes
				auth := v1.Group("/auth")
				{
					// Public routes (no auth required)
					auth.POST("/register", gatewayHandler.ProxyRequest)
					auth.POST("/login", gatewayHandler.ProxyRequest)
				}

				// Protected identity service routes
				protectedIdentity := v1.Group("")
				protectedIdentity.Use(middleware.AuthMiddleware(&cfg.JWT, logger))
				{
					users := protectedIdentity.Group("/users")
					{
						users.GET("/profile", gatewayHandler.ProxyRequest)
						users.PUT("/profile", gatewayHandler.ProxyRequest)
						users.PUT("/password", gatewayHandler.ProxyRequest)
					}

					addresses := protectedIdentity.Group("/addresses")
					{
						addresses.GET("", gatewayHandler.ProxyRequest)
						addresses.POST("", gatewayHandler.ProxyRequest)
						addresses.GET("/:id", gatewayHandler.ProxyRequest)
						addresses.PUT("/:id", gatewayHandler.ProxyRequest)
						addresses.DELETE("/:id", gatewayHandler.ProxyRequest)
						addresses.PUT("/:id/default", gatewayHandler.ProxyRequest)
					}
				}
			}
		}

	// Catch-all route for any unmatched paths
	router.NoRoute(gatewayHandler.ProxyRequest)

	return router
}

// InitializeServices registers all microservices from configuration
func InitializeServices(cfg *config.Config, serviceRegistry interface{}, logger *zap.Logger) error {
	// This would be implemented to register services from config
	// For now, services are registered in main.go
	return nil
}

