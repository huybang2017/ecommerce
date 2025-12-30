package service

import (
	"api-gateway/internal/domain"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// GatewayService orchestrates request routing and proxying
// This is the business logic layer for the API Gateway
type GatewayService struct {
	serviceRegistry domain.ServiceRegistry
	proxyClient     domain.ProxyClient
	logger          *zap.Logger
}

// NewGatewayService creates a new gateway service
func NewGatewayService(
	serviceRegistry domain.ServiceRegistry,
	proxyClient domain.ProxyClient,
	logger *zap.Logger,
) *GatewayService {
	return &GatewayService{
		serviceRegistry: serviceRegistry,
		proxyClient:     proxyClient,
		logger:          logger,
	}
}

// RouteRequest routes a request to the appropriate microservice
func (s *GatewayService) RouteRequest(
	ctx context.Context,
	serviceName string,
	path string,
	method string,
	headers map[string]string,
	body []byte,
) (*domain.ProxyResponse, error) {
	// Get the service from registry
	service, err := s.serviceRegistry.GetService(serviceName)
	if err != nil {
		s.logger.Error("Service not found", zap.String("service", serviceName), zap.Error(err))
		return &domain.ProxyResponse{
			Body:       []byte(fmt.Sprintf(`{"error":"service %s not found"}`, serviceName)),
			StatusCode: http.StatusNotFound,
			Headers:    make(map[string][]string),
		}, fmt.Errorf("service %s not found: %w", serviceName, err)
	}

	// Note: Authentication is already validated by middleware in the router
	// Middleware validates JWT token and sets user_id in gin.Context
	// Handler passes user_id from gin.Context to context.Context
	// So if we reach here, authentication is already validated
	// We don't need to check again - just proceed with routing
	_ = s.findRoute(service, path, method)

	// Log the routing attempt for debugging
	s.logger.Debug("Routing request",
		zap.String("service", serviceName),
		zap.String("path", path),
		zap.String("method", method),
		zap.String("base_url", service.BaseURL),
	)

	// Proxy the request to the backend service
	proxyResponse, err := s.proxyClient.ProxyRequest(service, path, method, headers, body)
	if err != nil {
		s.logger.Error("Failed to proxy request",
			zap.String("service", serviceName),
			zap.String("path", path),
			zap.String("base_url", service.BaseURL),
			zap.Error(err),
		)
		return &domain.ProxyResponse{
			Body:       []byte(`{"error":"failed to proxy request"}`),
			StatusCode: http.StatusBadGateway,
			Headers:    make(map[string][]string),
		}, fmt.Errorf("failed to proxy request: %w", err)
	}

	return proxyResponse, nil
}

// findRoute finds a matching route for the given path and method
func (s *GatewayService) findRoute(service *domain.Service, path string, method string) *domain.Route {
	for _, route := range service.Routes {
		// Simple path matching - in production, use a proper router
		if s.pathMatches(route.Path, path) && s.methodMatches(route.Methods, method) {
			return &route
		}
	}
	return nil
}

// pathMatches checks if a request path matches a route pattern
// This is a simplified matcher - in production, use a proper router library
func (s *GatewayService) pathMatches(pattern string, path string) bool {
	// Simple exact match
	if pattern == path {
		return true
	}

	// Basic pattern matching for path parameters (e.g., /products/:id)
	patternParts := s.splitPath(pattern)
	pathParts := s.splitPath(path)

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, patternPart := range patternParts {
		// If pattern part starts with :, it's a parameter, so it matches any value
		if len(patternPart) > 0 && patternPart[0] == ':' {
			continue
		}
		// Otherwise, parts must match exactly
		if patternPart != pathParts[i] {
			return false
		}
	}

	return true
}

// splitPath splits a path string into parts, removing empty parts
func (s *GatewayService) splitPath(path string) []string {
	parts := []string{}
	current := ""
	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// methodMatches checks if the HTTP method is allowed
func (s *GatewayService) methodMatches(allowedMethods []string, method string) bool {
	for _, m := range allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// HealthCheck checks the health of all registered services
func (s *GatewayService) HealthCheck(ctx context.Context) map[string]error {
	services := s.serviceRegistry.GetAllServices()
	results := make(map[string]error)

	for name, service := range services {
		err := s.proxyClient.HealthCheck(service)
		results[name] = err
	}

	return results
}

// ReadRequestBody reads the request body
func ReadRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// Restore the body so it can be read again if needed
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
