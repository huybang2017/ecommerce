package router

import (
	"api-gateway/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	cfg *config.Config,
	logger *zap.Logger,
	redisClient *redis.Client,
) *gin.Engine {
	// Use gin.New() instead of gin.Default() to avoid default middlewares
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// CRITICAL: Custom CORS middleware MUST be first
	router.Use(middleware.CORSMiddleware(&cfg.CORS, logger))

	// Skip logging OPTIONS requests (CORS preflight) to reduce noise
	router.Use(middleware.SkipOptionsLoggingMiddleware(logger))

	// Request logging middleware
	router.Use(middleware.RequestLoggingMiddleware(logger))
	router.Use(middleware.ErrorLoggingMiddleware(logger))

	// Rate limiting middleware
	router.Use(middleware.RateLimitMiddleware(&cfg.RateLimit, logger))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Swagger multi-service configuration
	router.GET("/swagger-config", func(c *gin.Context) {
		config := gin.H{
			"urls": []gin.H{
				{"url": "/swagger/services/product", "name": "Product Service"},
				{"url": "/swagger/services/order", "name": "Order Service"},
				{"url": "/swagger/services/identity", "name": "Identity Service"},
				{"url": "/swagger/services/search", "name": "Search Service"},
			},
		}
		c.JSON(http.StatusOK, config)
	})

	// Proxy swagger docs from microservices
	router.GET("/swagger/services/product", gatewayHandler.ProxySwagger("product_service"))
	router.GET("/swagger/services/order", gatewayHandler.ProxySwagger("order_service"))
	router.GET("/swagger/services/identity", gatewayHandler.ProxySwagger("identity_service"))
	router.GET("/swagger/services/search", gatewayHandler.ProxySwagger("search_service"))

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

				// Product Items (SKU) routes - Public
				products.GET("/:id/items", productHandler.GetProductItems)
				products.GET("/:id/items/:item_id", productHandler.GetProductItem)

				// Variation routes - Public (for UI selectors)
				products.GET("/:id/variations", productHandler.GetProductVariations)

				products.POST("", productHandler.CreateProduct) // Protected in handler

				// Protected routes (auth required)
				protected := products.Group("")
				protected.Use(middleware.AuthMiddleware(&cfg.JWT, logger), middleware.SessionMiddleware(logger, redisClient))
				{
					protected.PUT("/:id", productHandler.UpdateProduct)
					protected.PATCH("/:id", productHandler.UpdateProduct)
					protected.PATCH("/:id/inventory", productHandler.UpdateInventory)
					protected.DELETE("/:id", productHandler.DeleteProduct)

					// Product Items (SKU) - Protected operations
					protected.POST("/:id/items", productHandler.CreateProductItem)
					protected.PUT("/:id/items/:item_id", productHandler.UpdateProductItem)
					protected.DELETE("/:id/items/:item_id", productHandler.DeleteProductItem)
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

			// Cart routes (Order Service) - Protected routes (require authentication)
			cart := v1.Group("/cart")
			cart.Use(middleware.AuthMiddleware(&cfg.JWT, logger))
			{
				cart.GET("", cartHandler.GetCart)
				cart.DELETE("", cartHandler.ClearCart)
				cart.POST("/items", cartHandler.AddItem)
				cart.PUT("/items/:product_item_id", cartHandler.UpdateItem)
				cart.DELETE("/items/:product_item_id", cartHandler.RemoveItem)
				cart.POST("/items/:product_item_id/toggle", cartHandler.ToggleItemSelection) // toggle selection
				// New PATCH routes to set selection state (idempotent)
				cart.PATCH("/items/:product_item_id/selection", cartHandler.SetItemSelection) // set selection for single item
				cart.PATCH("/selection", cartHandler.SetAllSelection)                        // set selection for all items
				cart.PATCH("/shops/:shop_id/selection", cartHandler.SetShopSelection)         // set selection by shop

				cart.DELETE("/selected", cartHandler.ClearSelected) // clear selected items
				cart.POST("/validate", cartHandler.ValidateCart)   // validate cart before checkout
			}

			// Order routes (Order Service)
			orders := v1.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.ListOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.GET("/number/:order_number", orderHandler.GetOrderByNumber)
			}

			// Identity service routes - Auth
			auth := v1.Group("/auth")
			{
				// Public routes (no auth required)
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
				auth.POST("/refresh", authHandler.RefreshToken) // Refresh access token
			}

			// Logout requires auth to get user_id
			authProtected := v1.Group("/auth")
			authProtected.Use(middleware.AuthMiddleware(&cfg.JWT, logger), middleware.SessionMiddleware(logger, redisClient))
			{
				authProtected.POST("/logout", authHandler.Logout)
			}

			// Protected identity service routes
			protectedIdentity := v1.Group("")
			protectedIdentity.Use(middleware.AuthMiddleware(&cfg.JWT, logger), middleware.SessionMiddleware(logger, redisClient))
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

	// REMOVED: NoRoute catch-all prevents CORS middleware from working properly
	// If you need fallback routing, handle it in specific route groups
	// router.NoRoute(gatewayHandler.ProxyRequest)

	return router
}

// InitializeServices registers all microservices from configuration
func InitializeServices(cfg *config.Config, serviceRegistry interface{}, logger *zap.Logger) error {
	// This would be implemented to register services from config
	// For now, services are registered in main.go
	return nil
}
