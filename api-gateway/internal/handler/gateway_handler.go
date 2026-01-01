package handler

import (
	"api-gateway/internal/service"
	"context"
	"net/http"
	"strings"

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

// isCORSHeader checks if a header is a CORS-related header (case-insensitive)
func isCORSHeader(key string) bool {
	lower := strings.ToLower(key)
	return strings.HasPrefix(lower, "access-control-")
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
// @Summary Proxy request to microservice
// @Description Routes requests to appropriate backend microservices (Identity Service, Product Service, etc.)
// @Tags Gateway
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token for protected routes"
// @Success 200 {object} map[string]interface{} "Response from backend service"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/{path} [get]
// @Router /api/v1/{path} [post]
// @Router /api/v1/{path} [put]
// @Router /api/v1/{path} [patch]
// @Router /api/v1/{path} [delete]
func (h *GatewayHandler) ProxyRequest(c *gin.Context) {
	// FIX 1: CRITICAL - OPTIONS should never reach here (handled by CORS middleware)
	if c.Request.Method == "OPTIONS" {
		h.logger.Error("OPTIONS request reached ProxyRequest - CORS middleware may have failed!")
		c.AbortWithStatus(204)
		return
	}

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

	// CRITICAL: Forward user context headers to backend microservices
	// This allows backend services to know WHO is making the request without re-validating JWT
	// This is the BFF (Backend For Frontend) pattern
	if userID, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userID.(string); ok {
			headers["X-User-Id"] = userIDStr
			h.logger.Debug("Forwarding X-User-Id header", zap.String("user_id", userIDStr))
		}
	}

	if email, exists := c.Get("email"); exists {
		if emailStr, ok := email.(string); ok {
			headers["X-User-Email"] = emailStr
		}
	}

	if role, exists := c.Get("role"); exists {
		if roleStr, ok := role.(string); ok {
			headers["X-User-Role"] = roleStr
		}
	}

	// Get user_id from gin.Context (set by auth middleware) and add to context
	ctx := c.Request.Context()
	if userID, exists := c.Get("user_id"); exists {
		ctx = context.WithValue(ctx, "user_id", userID)
	}

	// Build path with query parameters
	fullPath := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		fullPath = fullPath + "?" + c.Request.URL.RawQuery
	}

	// Route the request
	proxyResponse, err := h.gatewayService.RouteRequest(
		ctx,
		serviceName,
		fullPath,
		c.Request.Method,
		headers,
		body,
	)

	if err != nil {
		statusCode := http.StatusBadGateway
		if proxyResponse != nil {
			statusCode = proxyResponse.StatusCode
		}

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
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// CRITICAL: Forward response headers from backend to client
	// EXCEPT CORS headers which are handled by Gateway middleware
	// This is essential for Set-Cookie headers in authentication
	h.logger.Info("Forwarding response headers",
		zap.Int("header_count", len(proxyResponse.Headers)),
		zap.Strings("header_keys", func() []string {
			keys := make([]string, 0, len(proxyResponse.Headers))
			for k := range proxyResponse.Headers {
				keys = append(keys, k)
			}
			return keys
		}()),
	)

	// FIX 2: Skip ALL CORS headers from backend (case-insensitive)
	for headerKey, headerValues := range proxyResponse.Headers {
		// Skip CORS headers - Gateway middleware handles them
		if isCORSHeader(headerKey) {
			h.logger.Debug("Skipping CORS header from backend",
				zap.String("key", headerKey),
			)
			continue
		}

		for _, headerValue := range headerValues {
			h.logger.Info("Setting response header",
				zap.String("key", headerKey),
				zap.String("value", headerValue),
			)
			c.Writer.Header().Add(headerKey, headerValue)
		}
	}

	// Get Content-Type from backend response headers
	contentType := "application/json"
	if ctValues, ok := proxyResponse.Headers["Content-Type"]; ok && len(ctValues) > 0 {
		contentType = ctValues[0]
	}

	// FIX 3: Verify CORS headers are present before sending response
	h.logger.Info("Final response headers before c.Data()",
		zap.String("Access-Control-Allow-Origin", c.Writer.Header().Get("Access-Control-Allow-Origin")),
		zap.String("Access-Control-Allow-Credentials", c.Writer.Header().Get("Access-Control-Allow-Credentials")),
		zap.Int("status_code", proxyResponse.StatusCode),
	)

	// Write response with backend's content type
	c.Data(proxyResponse.StatusCode, contentType, proxyResponse.Body)
}

// HealthCheck returns the health status of the gateway and all services
// @Summary Health check
// @Description Returns the health status of the API Gateway and all registered microservices
// @Tags Gateway
// @Produce json
// @Success 200 {object} map[string]interface{} "Gateway and services are healthy"
// @Failure 503 {object} map[string]interface{} "One or more services are unhealthy"
// @Router /health [get]
// @Router /api/gateway/health [get]
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
			"status":   "healthy",
			"gateway":  "ok",
			"services": healthStatus,
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"gateway":  "ok",
			"services": healthStatus,
		})
	}
}

// getServiceName maps request paths to service names
func (h *GatewayHandler) getServiceName(path string) string {
	// Simple path-based routing
	// IMPORTANT: Order matters! More specific paths must be checked first
	if strings.HasPrefix(path, "/api/v1/search") {
		return "search_service"
	}
	if strings.HasPrefix(path, "/api/v1/products/search") {
		// This is Product Service's search endpoint, not Search Service
		return "product_service"
	}
	if strings.HasPrefix(path, "/api/v1/products") {
		return "product_service"
	}
	if strings.HasPrefix(path, "/api/v1/categories") {
		return "product_service"
	}
	if strings.HasPrefix(path, "/api/v1/auth") {
		return "identity_service"
	}
	if strings.HasPrefix(path, "/api/v1/users") {
		return "identity_service"
	}
	if strings.HasPrefix(path, "/api/v1/addresses") {
		return "identity_service"
	}
	if strings.HasPrefix(path, "/api/v1/shops") { // THÊM MỚI - Shop routes
		return "identity_service"
	}
	if strings.HasPrefix(path, "/api/v1/cart") {
		return "order_service"
	}
	if strings.HasPrefix(path, "/api/v1/orders") {
		return "order_service"
	}
	// Default to product_service for now
	return "product_service"
}
