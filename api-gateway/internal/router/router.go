package router

import (
	"api-gateway/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// SetupRouter configures all API Gateway routes
func SetupRouter(
	gatewayHandler *handler.GatewayHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	addressHandler *handler.AddressHandler,
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	searchHandler *handler.SearchHandler,
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

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
					products.GET("", productHandler.ListProducts)
					products.GET("/:id", productHandler.GetProduct)
					products.GET("/search", productHandler.SearchProducts)
					products.POST("", productHandler.CreateProduct) // Protected in handler

					// Protected routes (auth required)
					protected := products.Group("")
					protected.Use(middleware.AuthMiddleware(&cfg.JWT, logger))
					{
						protected.PUT("/:id", productHandler.UpdateProduct)
						protected.PATCH("/:id", productHandler.UpdateProduct)
						protected.PATCH("/:id/inventory", productHandler.UpdateInventory)
						protected.DELETE("/:id", productHandler.DeleteProduct)
					}
				}

				// Category routes (Product Service)
				categories := v1.Group("/categories")
				{
					// Public routes (no auth required)
					categories.GET("", categoryHandler.ListCategories)
					categories.GET("/:id", categoryHandler.GetCategory)
					categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
					categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
					categories.GET("/:id/products", categoryHandler.GetCategoryProducts)
					categories.POST("", categoryHandler.CreateCategory)
					categories.PUT("/:id", categoryHandler.UpdateCategory)
					categories.DELETE("/:id", categoryHandler.DeleteCategory)
				}

				// Search routes (Search Service)
				search := v1.Group("/search")
				{
					search.GET("", searchHandler.SearchProducts)
				}

				// Identity service routes - Auth
				auth := v1.Group("/auth")
				{
					// Public routes (no auth required)
					auth.POST("/register", authHandler.Register)
					auth.POST("/login", authHandler.Login)
				}

				// Protected identity service routes
				protectedIdentity := v1.Group("")
				protectedIdentity.Use(middleware.AuthMiddleware(&cfg.JWT, logger))
				{
					users := protectedIdentity.Group("/users")
					{
						users.GET("/profile", userHandler.GetProfile)
						users.PUT("/profile", userHandler.UpdateProfile)
						users.PUT("/password", userHandler.ChangePassword)
					}

					addresses := protectedIdentity.Group("/addresses")
					{
						addresses.GET("", addressHandler.GetAddresses)
						addresses.POST("", addressHandler.CreateAddress)
						addresses.GET("/:id", addressHandler.GetAddress)
						addresses.PUT("/:id", addressHandler.UpdateAddress)
						addresses.DELETE("/:id", addressHandler.DeleteAddress)
						addresses.PUT("/:id/default", addressHandler.SetDefaultAddress)
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

