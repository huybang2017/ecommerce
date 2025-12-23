package repository

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
	"api-gateway/internal/domain"
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
) ([]byte, int, error) {
	// Build the full URL
	// Ensure base URL doesn't end with / and path starts with /
	baseURL := service.BaseURL
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
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// CRITICAL: Set ALL headers from map to request
	// This ensures Authorization header is always forwarded
	for key, value := range headers {
		if key == "" || value == "" {
			continue
		}
		req.Header.Set(key, value)
	}
	
	// CRITICAL: Double-check Authorization header is set
	// If it's in the headers map, ensure it's in the request
	if authVal, exists := headers["Authorization"]; exists && authVal != "" {
		// Force set it again to be absolutely sure
		req.Header.Set("Authorization", authVal)
		fmt.Printf("[PROXY] ✅ Set Authorization: %s...\n", authVal[:min(50, len(authVal))])
		
		// Verify it's actually in the request
		if finalAuth := req.Header.Get("Authorization"); finalAuth != "" {
			fmt.Printf("[PROXY] ✅ Verified Authorization in request\n")
		} else {
			fmt.Printf("[PROXY] ❌ ERROR: Authorization missing after setting!\n")
		}
	} else {
		fmt.Printf("[PROXY] ❌ ERROR: Authorization NOT in headers map! Keys: %v\n", getHeaderKeys(headers))
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
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
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
