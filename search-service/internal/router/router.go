package router

import (
	"os"
	"search-service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all API routes
// This is the transport layer - it defines the HTTP API surface
func SetupRouter(searchHandler *handler.SearchHandler) *gin.Engine {
	router := gin.Default()

	// Lightweight CORS middleware for development (enables Swagger "Try it out")
	devCORS := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-User-Id, X-User-Email, X-User-Role")
			c.Header("Access-Control-Allow-Credentials", "true")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		}
	}()
	// Enable dev CORS middleware by default. To disable set env ENABLE_DEV_CORS=false
	if os.Getenv("ENABLE_DEV_CORS") != "false" {
		router.Use(devCORS)
	}

	// Health check endpoint
	router.GET("/health", searchHandler.HealthCheck)

	// Swagger UI (serve generated docs)
	// Keep the UI on the service for dev convenience
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Expose OpenAPI spec at service-scoped path for Gateway to consume
	// e.g. GET /search/swagger/doc.json
	router.GET("/search/swagger/doc.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.File("docs/swagger.json")
	})
	// Preflight support (handled by devCORS middleware)
	router.OPTIONS("/search/swagger/doc.json", func(c *gin.Context) {
		c.Status(204)
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Search routes
		v1.GET("/search", searchHandler.SearchProducts)
	}

	return router
}
