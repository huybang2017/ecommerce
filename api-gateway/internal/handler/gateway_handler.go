package handler

import (
	"context"
	"net/http"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getHeaderKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GatewayHandler handles HTTP requests for the API Gateway
type GatewayHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewGatewayHandler creates a new gateway handler
func NewGatewayHandler(gatewayService *service.GatewayService, logger *zap.Logger) *GatewayHandler {
	return &GatewayHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// ProxyRequest proxies a request to the appropriate microservice
// This is the main handler that routes requests to backend services
func (h *GatewayHandler) ProxyRequest(c *gin.Context) {
	// DEBUG: Log incoming request immediately
	_, hasAuthInContext := c.Get("auth_header")
	h.logger.Info("ProxyRequest called",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("auth_header_in_request", c.Request.Header.Get("Authorization")),
		zap.Bool("auth_in_context", hasAuthInContext),
	)
	
	// Extract service name from path
	serviceName := h.getServiceName(c.Request.URL.Path)

	// Read request body
	body, err := service.ReadRequestBody(c.Request)
	if err != nil {
		h.logger.Error("Failed to read request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Collect headers - CRITICAL: Always include Authorization header
	headers := make(map[string]string)
	
	// FIRST: Copy ALL headers from request (including Authorization)
	// This ensures we don't miss any headers
	for key, values := range c.Request.Header {
		// Skip hop-by-hop headers that shouldn't be forwarded
		if key == "Connection" || key == "Keep-Alive" || key == "Transfer-Encoding" || key == "Upgrade" {
			continue
		}
		// Copy all other headers including Authorization
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	
	// CRITICAL: Ensure Authorization header is present
	// Priority 1: Get from context (preserved by middleware)
	var authHeader string
	if preservedAuth, exists := c.Get("auth_header"); exists {
		if authStr, ok := preservedAuth.(string); ok && authStr != "" {
			authHeader = authStr
			// Override with preserved header from middleware
			headers["Authorization"] = authHeader
			h.logger.Debug("Got Authorization from context", zap.String("header_preview", authStr[:min(30, len(authStr))]))
		}
	}
	
	// Priority 2: Get from Request.Header if not in context
	if authHeader == "" {
		authHeader = c.Request.Header.Get("Authorization")
		if authHeader != "" {
			headers["Authorization"] = authHeader
			h.logger.Debug("Got Authorization from Request.Header", zap.String("header_preview", authHeader[:min(30, len(authHeader))]))
		}
	}
	
	// Final check: Log if Authorization is missing
	if headers["Authorization"] == "" {
		h.logger.Warn("No Authorization header found in handler", zap.Strings("available_headers", getHeaderKeys(headers)))
	} else {
		h.logger.Debug("Authorization header ready for forwarding", zap.String("header_preview", headers["Authorization"][:min(30, len(headers["Authorization"]))]))
	}

	// Get user_id from gin.Context (set by auth middleware) and add to context
	ctx := c.Request.Context()
	if userID, exists := c.Get("user_id"); exists {
		ctx = context.WithValue(ctx, "user_id", userID)
	}

	// Route the request
	responseBody, statusCode, err := h.gatewayService.RouteRequest(
		ctx,
		serviceName,
		c.Request.URL.Path,
		c.Request.Method,
		headers,
		body,
	)

	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Failed to route request",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("service", serviceName),
		)
		c.JSON(statusCode, gin.H{
			"error": "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// Set response headers
	c.Header("Content-Type", "application/json")

	// Write response
	c.Data(statusCode, "application/json", responseBody)
}

// HealthCheck returns the health status of the gateway and all services
func (h *GatewayHandler) HealthCheck(c *gin.Context) {
	healthStatus := h.gatewayService.HealthCheck(c.Request.Context())

	allHealthy := true
	for serviceName, err := range healthStatus {
		if err != nil {
			allHealthy = false
			h.logger.Warn("Service unhealthy", zap.String("service", serviceName), zap.Error(err))
		}
	}

	if allHealthy {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"gateway": "ok",
			"services": healthStatus,
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "degraded",
			"gateway": "ok",
			"services": healthStatus,
		})
	}
}

// getServiceName maps request paths to service names
func (h *GatewayHandler) getServiceName(path string) string {
	// Simple path-based routing
	if len(path) >= 12 && path[:12] == "/api/v1/prod" {
		return "product_service"
	}
	if len(path) >= 15 && path[:15] == "/api/v1/categor" {
		return "product_service"
	}
	if len(path) >= 12 && path[:12] == "/api/v1/auth" {
		return "identity_service"
	}
	if len(path) >= 12 && path[:12] == "/api/v1/user" {
		return "identity_service"
	}
	if len(path) >= 15 && path[:15] == "/api/v1/address" {
		return "identity_service"
	}
	// Default to product_service for now
	return "product_service"
}
