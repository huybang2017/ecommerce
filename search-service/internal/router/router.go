package router

import (
	"search-service/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all API routes
// This is the transport layer - it defines the HTTP API surface
func SetupRouter(searchHandler *handler.SearchHandler) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", searchHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Search routes
		v1.GET("/search", searchHandler.SearchProducts)
	}

	return router
}


