package router

import (
	"product-service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "product-service/docs" // Import generated docs
)

// SetupRouter configures all API routes
// This is the transport layer - it defines the HTTP API surface
func SetupRouter(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler) *gin.Engine {
	router := gin.Default()

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.ListProducts)          // List products with pagination and filters
			products.POST("", productHandler.CreateProduct)
			products.GET("/search", productHandler.SearchProducts) // Search (must be before /:id)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.PATCH("/:id/inventory", productHandler.UpdateInventory)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug) // Must be before /:id
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
			categories.GET("/:id/products", productHandler.GetProductsByCategory) // Products by category
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}
	}

	return router
}

