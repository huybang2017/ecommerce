package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
	Services  ServicesConfig
	Logging   LoggingConfig
	Redis     RedisConfig
}

type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	Burst             int
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

type ServiceConfig struct {
	BaseURL         string
	Timeout         time.Duration
	HealthCheckPath string
	Routes          []RouteConfig
}

type RouteConfig struct {
	Path        string
	Methods     []string
	RequireAuth bool
}

type ServicesConfig map[string]ServiceConfig

type LoggingConfig struct {
	Level            string
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v. Using defaults and environment variables.", err)
	}

	config := &Config{}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	services := make(ServicesConfig)

	serviceKeys := []string{"product_service", "identity_service", "search_service", "order_service"}
	for _, serviceKey := range serviceKeys {
		servicePath := fmt.Sprintf("services.%s", serviceKey)
		envVarName := fmt.Sprintf("SERVICES_%s_BASE_URL", strings.ToUpper(strings.ReplaceAll(serviceKey, "_", "_")))
		baseURL := os.Getenv(envVarName)
		serviceConfig := ServiceConfig{
			BaseURL:         baseURL,
			Timeout:         viper.GetDuration(fmt.Sprintf("%s.timeout", servicePath)),
			HealthCheckPath: viper.GetString(fmt.Sprintf("%s.health_check_path", servicePath)),
		}
		if baseURL == "" {
			serviceConfig.BaseURL = viper.GetString(fmt.Sprintf("%s.base_url", servicePath))
		}

		routesPath := fmt.Sprintf("%s.routes", servicePath)
		if viper.IsSet(routesPath) {
			var routes []RouteConfig
			if err := viper.UnmarshalKey(routesPath, &routes); err == nil {
				serviceConfig.Routes = routes
			}
		}

		if serviceConfig.BaseURL != "" {
			services[serviceKey] = serviceConfig
		}
	}

	if len(services) > 0 {
		config.Services = services
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration", "24h")
	viper.SetDefault("jwt.issuer", "api-gateway")

	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests_per_minute", 100)
	viper.SetDefault("rate_limit.burst", 20)

	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:5173"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie", "Set-Cookie"})
	viper.SetDefault("cors.expose_headers", []string{"Set-Cookie"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", "12h")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	viper.SetDefault("services.product_service.base_url", "http://localhost:8080")
	viper.SetDefault("services.product_service.timeout", "30s")
	viper.SetDefault("services.product_service.health_check_path", "/health")

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")
	viper.SetDefault("logging.output_paths", []string{"stdout"})
	viper.SetDefault("logging.error_output_paths", []string{"stderr"})
}

func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
