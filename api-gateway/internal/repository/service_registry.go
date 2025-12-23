package repository

import (
	"fmt"
	"api-gateway/internal/domain"
	"sync"
)

// serviceRegistry implements the ServiceRegistry interface
// This is an in-memory service registry
// In production, you might use Consul, Eureka, or Kubernetes service discovery
type serviceRegistry struct {
	services map[string]*domain.Service
	mu       sync.RWMutex
}

// NewServiceRegistry creates a new in-memory service registry
func NewServiceRegistry() domain.ServiceRegistry {
	return &serviceRegistry{
		services: make(map[string]*domain.Service),
	}
}

// GetService retrieves a service by name
func (r *serviceRegistry) GetService(name string) (*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// GetAllServices returns all registered services
func (r *serviceRegistry) GetAllServices() map[string]*domain.Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]*domain.Service)
	for k, v := range r.services {
		result[k] = v
	}
	return result
}

// RegisterService registers a new service
func (r *serviceRegistry) RegisterService(service *domain.Service) error {
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}
	if service.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.services[service.Name] = service
	return nil
}

