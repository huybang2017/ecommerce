package handler

import (
	"api-gateway/internal/service"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

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
// This is the main handler that routes requests to backend microservices
func (h *GatewayHandler) ProxyRequest(c *gin.Context) {
	// FIX 1: CRITICAL - OPTIONS should never reach here (handled by CORS middleware)
	if c.Request.Method == "OPTIONS" {
		h.logger.Error("OPTIONS request reached ProxyRequest - CORS middleware may have failed!")
		c.AbortWithStatus(204)
		return
	}

	// Minimal logging: method and path
	h.logger.Debug("ProxyRequest", zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))

	// Extract service name from path
	serviceName := h.getServiceName(c.Request.URL.Path)

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
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

	// Forward user context headers to backend microservices (injected by gateway)
	if userID, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userID.(string); ok {
			headers["X-User-Id"] = userIDStr
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

	// Forward response headers from backend to client (except CORS)

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

// SpecProxy fetches a service's swagger.json server-side, rewrites the
// `servers` entry to point to the gateway, and returns the modified spec.
func (h *GatewayHandler) SpecProxy(c *gin.Context) {
	serviceParam := c.Param("service")
	// Fallback: some router edge-cases may not populate c.Param correctly.
	// Parse from URL path if param is empty.
	if serviceParam == "" {
		raw := c.Request.URL.Path
		h.logger.Debug("SpecProxy: c.Param empty, parsing from URL", zap.String("path", raw))
		// trim prefix
		if strings.HasPrefix(raw, "/specs/") {
			raw = strings.TrimPrefix(raw, "/specs/")
		}
		// remove leading/trailing slashes
		raw = strings.Trim(raw, "/ ")
		// if it contains query params, strip
		if idx := strings.Index(raw, "?"); idx != -1 {
			raw = raw[:idx]
		}
		serviceParam = strings.TrimSuffix(raw, ".json")
	}
	if serviceParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing service parameter"})
		return
	}

	// Trim common suffixes if the route captured them
	serviceParam = strings.TrimSuffix(serviceParam, ".json")
	serviceParam = strings.Trim(serviceParam, "/ ")

	// Map short name (e.g. product) -> service key (product_service)
	svcKey := strings.ReplaceAll(serviceParam, "-", "_") + "_service"

	sr := h.gatewayService.GetServiceRegistry()
	svc, err := sr.GetService(svcKey)
	if err != nil {
		// Try exact key as fallback (serviceParam might already be full key)
		svc, err = sr.GetService(serviceParam)
		if err != nil {
			h.logger.Warn("SpecProxy: service not found", zap.String("service_param", serviceParam), zap.String("svc_key", svcKey))
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found", "service": serviceParam})
			return
		}
	}

	// Construct remote spec URL (service is expected to expose /<short>/swagger/doc.json)
	remoteURL := strings.TrimRight(svc.BaseURL, "/") + "/" + serviceParam + "/swagger/doc.json"
	h.logger.Debug("SpecProxy fetching remote spec", zap.String("remote_url", remoteURL))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(remoteURL)
	if err != nil {
		h.logger.Error("SpecProxy failed to fetch remote spec", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch remote spec", "details": err.Error(), "remote_url": remoteURL})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.logger.Warn("SpecProxy remote spec returned non-200", zap.Int("status", resp.StatusCode), zap.String("remote_url", remoteURL))
		c.JSON(http.StatusBadGateway, gin.H{"error": "remote spec returned non-200", "status": resp.StatusCode, "remote_url": remoteURL})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("SpecProxy failed reading remote body", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read remote spec"})
		return
	}

	// Try to unmarshal JSON and rewrite servers
	var doc map[string]interface{}
	if err := json.Unmarshal(body, &doc); err != nil {
		// If not JSON, just proxy raw body
		h.logger.Warn("SpecProxy: remote spec not JSON, proxying raw body")
		c.Data(http.StatusOK, "application/json", body)
		return
	}

	// Build gateway base URL (scheme + host) and set /api/v1 as base for runtime
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	gatewayBase := scheme + "://" + c.Request.Host + "/api/v1"

	// If OpenAPI v3 (has `openapi` or `servers`), rewrite `servers`.
	if _, ok := doc["openapi"]; ok {
		servers := []map[string]string{{"url": gatewayBase}}
		doc["servers"] = servers
	}

	// If Swagger 2.0 (has `swagger` == "2.0"), rewrite `host` and `basePath`.
	if v, ok := doc["swagger"]; ok {
		if s, ok2 := v.(string); ok2 && s == "2.0" {
			// For swagger2, set host and basePath so UI builds URLs like scheme://host{basePath}{path}
			doc["host"] = c.Request.Host
			doc["basePath"] = "/api/v1"
			// also remove servers if present
			delete(doc, "servers")
		}
	}

	out, err := json.Marshal(doc)
	if err != nil {
		h.logger.Error("SpecProxy failed to marshal modified spec", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate spec"})
		return
	}

	c.Data(http.StatusOK, "application/json", out)
}

// SetSessionCookie allows the Swagger UI to request the gateway set an httpOnly session cookie.
// This is intended for local/dev use to let the browser send the session cookie while keeping it httpOnly.
func (h *GatewayHandler) SetSessionCookie(c *gin.Context) {
	var payload struct {
		SessionID string `json:"session_id"`
		MaxAge    int    `json:"max_age"`
	}

	if err := c.BindJSON(&payload); err != nil {
		h.logger.Warn("SetSessionCookie: invalid payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Security: ensure this endpoint is only usable from the same origin (Swagger UI served by gateway)
	origin := c.Request.Header.Get("Origin")
	if origin != "" && !strings.Contains(origin, c.Request.Host) {
		h.logger.Warn("SetSessionCookie: origin mismatch", zap.String("origin", origin), zap.String("host", c.Request.Host))
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// If SessionID empty => clear cookie
	if payload.SessionID == "" {
		// set cookie with max-age 0 to delete
		c.SetCookie("session_id", "", 0, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"status": "cleared"})
		return
	}

	// Use provided MaxAge if positive, otherwise default (e.g., 3600 seconds)
	maxAge := payload.MaxAge
	if maxAge <= 0 {
		maxAge = 3600
	}

	// Set httpOnly cookie so browser will send it but JS cannot read it
	// For local dev, Secure=false; in production set Secure=true and SameSite as needed
	c.SetCookie("session_id", payload.SessionID, maxAge, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Note: Gateway does not proxy or merge service OpenAPI specs; it exposes /swagger-config which lists service spec URLs.

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
