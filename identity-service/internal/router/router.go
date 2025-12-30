package router

import (
	"identity-service/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all API routes
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	addressHandler *handler.AddressHandler,
	shopHandler *handler.ShopHandler,
	authMiddleware gin.HandlerFunc,
) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken) // Refresh access token
			auth.POST("/logout", authHandler.Logout)        // Logout (will need middleware for user_id)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(authMiddleware)
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.PUT("/password", userHandler.ChangePassword)
			}

			// Address routes
			addresses := protected.Group("/addresses")
			{
				addresses.GET("", addressHandler.GetAddresses)
				addresses.POST("", addressHandler.CreateAddress)
				addresses.GET("/:id", addressHandler.GetAddress)
				addresses.PUT("/:id", addressHandler.UpdateAddress)
				addresses.DELETE("/:id", addressHandler.DeleteAddress)
				addresses.PUT("/:id/default", addressHandler.SetDefaultAddress)
			}
		}

		// Shop routes
		shops := v1.Group("/shops")
		{
			// Public routes
			shops.GET("", shopHandler.ListShops)   // List all shops
			shops.GET("/:id", shopHandler.GetShop) // Get shop by ID
		}

		// Protected shop routes
		protectedShops := v1.Group("/shops")
		protectedShops.Use(authMiddleware)
		{
			protectedShops.POST("", shopHandler.CreateShop)                 // Create shop (SELLER only)
			protectedShops.GET("/my-shop", shopHandler.GetMyShop)           // Get my shop
			protectedShops.PUT("/:id", shopHandler.UpdateShop)              // Update shop (owner or ADMIN)
			protectedShops.DELETE("/:id", shopHandler.DeleteShop)           // Delete shop (ADMIN only)
			protectedShops.PUT("/:id/status", shopHandler.UpdateShopStatus) // Update status (ADMIN only)
		}
	}

	return router
}
