package router

import (
	"order-service/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all API routes
// This is the transport layer - it defines the HTTP API surface
func SetupRouter(cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) *gin.Engine {
	router := gin.Default()

	// CORS middleware - Allow frontend to access the API
	// Use Default() which handles OPTIONS automatically
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Explicitly handle OPTIONS for all routes (CORS preflight)
	// This must be registered before other routes
	router.OPTIONS("/api/v1/cart", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200")
		c.Status(204) // No Content
	})
	router.OPTIONS("/api/v1/cart/items", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200")
		c.Status(204)
	})
	router.OPTIONS("/api/v1/cart/items/:product_id", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200")
		c.Status(204)
	})

	// Health check endpoint
	router.GET("/health", cartHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Cart routes
		cart := v1.Group("/cart")
		{
			cart.GET("", cartHandler.GetCart)                    // Get cart
			cart.DELETE("", cartHandler.ClearCart)               // Clear cart
			cart.POST("/items", cartHandler.AddItem)             // Add item to cart
			cart.PUT("/items/:product_id", cartHandler.UpdateItem)   // Update item quantity
			cart.DELETE("/items/:product_id", cartHandler.RemoveItem) // Remove item from cart
		}

		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)                      // Create order from cart
			orders.GET("", orderHandler.ListOrders)                        // List orders
			orders.GET("/:id", orderHandler.GetOrder)                      // Get order by ID
			orders.GET("/number/:order_number", orderHandler.GetOrderByOrderNumber) // Get order by order number
		}
	}

	return router
}

