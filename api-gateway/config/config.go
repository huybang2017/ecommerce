package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the API Gateway
type Config struct {
	Server    ServerConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
	Services  ServicesConfig
	Logging   LoggingConfig
	Redis     RedisConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	Burst             int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// ServiceConfig holds configuration for a single microservice
type ServiceConfig struct {
	BaseURL         string
	Timeout         time.Duration
	HealthCheckPath string
	Routes          []RouteConfig
}

// RouteConfig defines a route pattern for a service
type RouteConfig struct {
	Path        string
	Methods     []string
	RequireAuth bool
}

// ServicesConfig holds configuration for all microservices
type ServicesConfig map[string]ServiceConfig

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level            string
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// LoadConfig reads configuration from config.yaml and environment variables
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("")

	// Set defaults
	setDefaults()

	// Read config file (optional - env vars will override)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v. Using defaults and environment variables.", err)
	}

	config := &Config{}

	// Unmarshal configuration into struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Fix: Manually unmarshal ServicesConfig because viper has issues with nested maps
	// Read directly from viper and construct ServiceConfig manually
	services := make(ServicesConfig)

	// Get all service keys
	serviceKeys := []string{"product_service", "identity_service", "search_service", "order_service"}
	for _, serviceKey := range serviceKeys {
		servicePath := fmt.Sprintf("services.%s", serviceKey)

		// Check for environment variable override first (e.g., SERVICES_PRODUCT_SERVICE_BASE_URL)
		envVarName := fmt.Sprintf("SERVICES_%s_BASE_URL", strings.ToUpper(strings.ReplaceAll(serviceKey, "_", "_")))
		baseURL := os.Getenv(envVarName)

		serviceConfig := ServiceConfig{
			BaseURL:         baseURL, // Use env var if set
			Timeout:         viper.GetDuration(fmt.Sprintf("%s.timeout", servicePath)),
			HealthCheckPath: viper.GetString(fmt.Sprintf("%s.health_check_path", servicePath)),
		}

		// If no env var, use config file value
		if baseURL == "" {
			serviceConfig.BaseURL = viper.GetString(fmt.Sprintf("%s.base_url", servicePath))
		}

		// Unmarshal routes separately
		routesPath := fmt.Sprintf("%s.routes", servicePath)
		if viper.IsSet(routesPath) {
			var routes []RouteConfig
			if err := viper.UnmarshalKey(routesPath, &routes); err == nil {
				serviceConfig.Routes = routes
			}
		}

		// Only add service if we have a base URL
		if serviceConfig.BaseURL != "" {
			services[serviceKey] = serviceConfig
		}
	}

	// Override Services with manually constructed values
	if len(services) > 0 {
		config.Services = services
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration", "24h")
	viper.SetDefault("jwt.issuer", "api-gateway")

	// Rate limit defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_minute", 100)
	viper.SetDefault("rate_limit.burst", 20)

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:5173"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie", "Set-Cookie"})
	viper.SetDefault("cors.expose_headers", []string{"Set-Cookie"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", "12h")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// Services defaults
	// Note: In Docker, use service name. For local dev, use localhost
	viper.SetDefault("services.product_service.base_url", "http://localhost:8080")
	viper.SetDefault("services.product_service.timeout", "30s")
	viper.SetDefault("services.product_service.health_check_path", "/health")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")
	viper.SetDefault("logging.output_paths", []string{"stdout"})
	viper.SetDefault("logging.error_output_paths", []string{"stderr"})
}

// GetAddress returns the Redis address
func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
