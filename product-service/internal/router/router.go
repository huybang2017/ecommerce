package router

import (
	"fmt"
	"log"
	"os"
	"product-service/docs"
	"product-service/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		fmt.Fprintf(os.Stderr, "📥📥📥 REQUEST RECEIVED: %s %s\n", method, path)
		log.Printf("📥 REQUEST RECEIVED: %s %s", method, path)
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Fprintf(os.Stderr, "📤📤📤 RESPONSE: %s %s - Status: %d - Latency: %v\n", method, path, status, latency)
		log.Printf("📤 RESPONSE: %s %s - Status: %d - Latency: %v", method, path, status, latency)
	}
}

func SetupRouter(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler, skuHandler *handler.SKUHandler, attrHandler *handler.AttributeHandler, stockHandler *handler.StockHandler, variationHandler *handler.VariationHandler) *gin.Engine {
	router := gin.Default()
	router.Use(RequestLogger())
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/product/swagger/doc.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.File("docs/swagger.json")
	})
	router.OPTIONS("/product/swagger/doc.json", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Status(204)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	v1 := router.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("", productHandler.ListProducts)
			products.POST("", productHandler.CreateProduct)
			products.GET("/search", productHandler.SearchProducts)

			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.PATCH("/:id/inventory", productHandler.UpdateInventory)

			products.GET("/:id/items", skuHandler.GetProductItems)
			products.POST("/:id/items", skuHandler.CreateProductItem)
			products.GET("/:id/items/:item_id", skuHandler.GetProductItem)
			products.PUT("/:id/items/:item_id", skuHandler.UpdateProductItem)
			products.DELETE("/:id/items/:item_id", skuHandler.DeleteProductItem)
			products.GET("/:id/variations", variationHandler.GetProductVariations)
			products.POST("/:id/attributes", attrHandler.SetProductAttributes)
			products.GET("/:id/attributes", attrHandler.GetProductAttributes)
		}
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.PATCH("/:id/active", categoryHandler.PatchCategoryActive)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
			categories.POST("/:id/children", categoryHandler.CreateChildCategory)
			categories.GET("/:id/products", productHandler.GetProductsByCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.POST("/:id/attributes", attrHandler.CreateCategoryAttribute)
			categories.GET("/:id/attributes", attrHandler.GetCategoryAttributes)
			categories.DELETE("/:id/attributes/:attr_id", attrHandler.DeleteCategoryAttribute)

			adminCategories := categories.Group("/admin")
			{
				adminCategories.GET("", categoryHandler.GetAdminCategories)
				adminCategories.GET("/parents", categoryHandler.GetAdminCategoryParent)
				adminCategories.POST("/parent", categoryHandler.CreateParentCategory)
			}
		}

		v1.GET("/product-items/batch", skuHandler.GetProductItemsBatch)
		v1.GET("/product-items/:id", skuHandler.GetProductItemBySKU)

		productItems := v1.Group("/product-items")
		{
			productItems.GET("/:id/stock", stockHandler.GetStock)
			productItems.PUT("/:id/stock", stockHandler.UpdateStock)
			productItems.POST("/check-stock", stockHandler.CheckStock)
			productItems.POST("/reserve-stock", stockHandler.ReserveStock)
			productItems.POST("/deduct-stock", stockHandler.DeductStock)
			productItems.POST("/release-stock", stockHandler.ReleaseStock)
		}
	}
	return router
}
