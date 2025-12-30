package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the Order Service
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Kafka         KafkaConfig
	Logging       LoggingConfig
	ProductService ProductServiceConfig
}

// ProductServiceConfig holds Product Service client configuration
type ProductServiceConfig struct {
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers          []string      `mapstructure:"brokers"`
	TopicOrderCreated string       `mapstructure:"topic_order_created"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	RequiredAcks     int           `mapstructure:"required_acks"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds PostgreSQL connection configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level            string
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
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

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8083)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5433)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "order_service")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic_order_created", "order_created")
	viper.SetDefault("kafka.write_timeout", "10s")
	viper.SetDefault("kafka.read_timeout", "10s")
	viper.SetDefault("kafka.required_acks", 1)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")
	viper.SetDefault("logging.output_paths", []string{"stdout"})
	viper.SetDefault("logging.error_output_paths", []string{"stderr"})

	// Product Service defaults
	viper.SetDefault("product_service.base_url", "http://localhost:8000")
	viper.SetDefault("product_service.timeout", "10s")
}

// GetDSN returns the PostgreSQL Data Source Name
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetAddress returns the Redis address
func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}


