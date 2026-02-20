package repository

import (
	"api-gateway/internal/domain"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// proxyClient implements the ProxyClient interface
// This handles HTTP proxying to backend microservices
type proxyClient struct {
	httpClient *http.Client
}

// NewProxyClient creates a new HTTP proxy client
func NewProxyClient(timeout time.Duration) domain.ProxyClient {
	return &proxyClient{
		httpClient: &http.Client{
			Timeout: timeout,
			// Don't follow redirects automatically
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// ProxyRequest proxies an HTTP request to a backend service
func (p *proxyClient) ProxyRequest(
	service *domain.Service,
	path string,
	method string,
	headers map[string]string,
	body []byte,
) (*domain.ProxyResponse, error) {
	// Build the full URL
	// Ensure base URL doesn't end with / and path starts with /
	baseURL := service.BaseURL
	// Safety net: if baseURL still points to a docker hostname, override to localhost for local dev
	if strings.Contains(baseURL, "product-service") {
		baseURL = "http://localhost:8080"
	}
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	url := baseURL + path

	// Create the request
	var req *http.Request
	var err error

	if body != nil && len(body) > 0 {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set ALL headers from map to request
	for key, value := range headers {
		if key == "" || value == "" {
			continue
		}
		req.Header.Set(key, value)
	}

	// Ensure Authorization header is explicitly set from the headers map
	if authVal, exists := headers["Authorization"]; exists && authVal != "" {
		req.Header.Set("Authorization", authVal)
	}

	// Set content type if body exists
	if body != nil && len(body) > 0 {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Execute the request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Capture ALL response headers, especially Set-Cookie
	responseHeaders := make(map[string][]string)
	for key, values := range resp.Header {
		responseHeaders[key] = values
	}

	return &domain.ProxyResponse{
		Body:       respBody,
		StatusCode: resp.StatusCode,
		Headers:    responseHeaders,
	}, nil
}

// HealthCheck checks if a service is healthy
func (p *proxyClient) HealthCheck(service *domain.Service) error {
	url := service.BaseURL + service.HealthCheckPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service unhealthy: status code %d", resp.StatusCode)
	}

	return nil
}