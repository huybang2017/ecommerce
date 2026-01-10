package router

import (
	"order-service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all API routes
// This is the transport layer - it defines the HTTP API surface
// NOTE: CORS is handled by API Gateway - this service should only receive internal requests
func SetupRouter(cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler) *gin.Engine {
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", cartHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Cart routes
		cart := v1.Group("/cart")
		{
			cart.GET("", cartHandler.GetCart)                              // Get cart
			cart.DELETE("", cartHandler.ClearCart)                         // Clear cart
			cart.POST("/items", cartHandler.AddItem)                       // Add item to cart
			cart.PUT("/items/:product_item_id", cartHandler.UpdateItem)    // Update item quantity
			cart.DELETE("/items/:product_item_id", cartHandler.RemoveItem) // Remove item from cart
		}

		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)                               // Create order from cart
			orders.GET("", orderHandler.ListOrders)                                 // List orders
			orders.GET("/:id", orderHandler.GetOrder)                               // Get order by ID
			orders.GET("/number/:order_number", orderHandler.GetOrderByOrderNumber) // Get order by order number
		}
	}

	return router
}
