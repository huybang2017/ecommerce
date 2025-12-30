package domain

// Service represents a backend microservice
// This is the domain model for service routing
type Service struct {
	Name            string
	BaseURL         string
	HealthCheckPath string
	Routes          []Route
}

// Route represents a route pattern for a service
type Route struct {
	Path        string
	Methods     []string
	RequireAuth bool
}

// ServiceRegistry defines the interface for service discovery
// This abstraction allows different service discovery mechanisms
type ServiceRegistry interface {
	GetService(name string) (*Service, error)
	GetAllServices() map[string]*Service
	RegisterService(service *Service) error
}

// ProxyResponse contains the full response from a proxied request
type ProxyResponse struct {
	Body       []byte
	StatusCode int
	Headers    map[string][]string
}

// ProxyClient defines the interface for proxying requests to services
// This abstraction allows different proxy implementations
type ProxyClient interface {
	ProxyRequest(service *Service, path string, method string, headers map[string]string, body []byte) (*ProxyResponse, error)
	HealthCheck(service *Service) error
}
