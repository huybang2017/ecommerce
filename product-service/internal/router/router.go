package router

import (
	"fmt"
	"log"
	"os"
	"product-service/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger middleware logs all incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Log request
		fmt.Fprintf(os.Stderr, "📥📥📥 REQUEST RECEIVED: %s %s\n", method, path)
		log.Printf("📥 REQUEST RECEIVED: %s %s", method, path)

		// Process request
		c.Next()

		// Log response
		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Fprintf(os.Stderr, "📤📤📤 RESPONSE: %s %s - Status: %d - Latency: %v\n", method, path, status, latency)
		log.Printf("📤 RESPONSE: %s %s - Status: %d - Latency: %v", method, path, status, latency)
	}
}

// SetupRouter configures all API routes
// This is the transport layer - it defines the HTTP API surface
func SetupRouter(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler, skuHandler *handler.SKUHandler, attrHandler *handler.AttributeHandler, stockHandler *handler.StockHandler) *gin.Engine {
	router := gin.Default()

	// Add request logging middleware
	router.Use(RequestLogger())

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
			products.GET("", productHandler.ListProducts) // List products with pagination and filters
			products.POST("", productHandler.CreateProduct)
			products.GET("/search", productHandler.SearchProducts) // Search (must be before /:id)

			// Product detail routes - MUST be first (before nested routes)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.PATCH("/:id/inventory", productHandler.UpdateInventory)

			// SKU routes (Product Items) - Use /:id/items (nested under product)
			products.GET("/:id/items", skuHandler.GetProductItems)               // List all SKUs for a product
			products.POST("/:id/items", skuHandler.CreateProductItem)            // Create new SKU
			products.GET("/:id/items/:item_id", skuHandler.GetProductItem)       // Get specific SKU
			products.PUT("/:id/items/:item_id", skuHandler.UpdateProductItem)    // Update SKU
			products.DELETE("/:id/items/:item_id", skuHandler.DeleteProductItem) // Delete SKU

			// Product attributes (EAV) - Use /:id/attributes
			products.POST("/:id/attributes", attrHandler.SetProductAttributes)
			products.GET("/:id/attributes", attrHandler.GetProductAttributes)
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

			// Category attributes (EAV) - Use /:id/attributes to avoid conflict
			categories.POST("/:id/attributes", attrHandler.CreateCategoryAttribute)
			categories.GET("/:id/attributes", attrHandler.GetCategoryAttributes)
			categories.DELETE("/:id/attributes/:attr_id", attrHandler.DeleteCategoryAttribute)
		}

		// Product item by SKU code (standalone route)
		v1.GET("/product-items/:id", skuHandler.GetProductItemBySKU)

		// Stock management routes
		productItems := v1.Group("/product-items")
		{
			productItems.GET("/:id/stock", stockHandler.GetStock)          // Get stock
			productItems.PUT("/:id/stock", stockHandler.UpdateStock)       // Update stock (shop owner)
			productItems.POST("/check-stock", stockHandler.CheckStock)     // Check stock availability
			productItems.POST("/reserve-stock", stockHandler.ReserveStock) // Reserve stock (checkout)
			productItems.POST("/deduct-stock", stockHandler.DeductStock)   // Deduct stock (payment confirmed)
			productItems.POST("/release-stock", stockHandler.ReleaseStock) // Release reservation (cancel/failed)
		}
	}

	return router
}
