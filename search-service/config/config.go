package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
// This is the single source of truth for configuration
type Config struct {
	Server        ServerConfig
	Kafka         KafkaConfig
	Elasticsearch ElasticsearchConfig
	Logging       LoggingConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// KafkaConfig holds Kafka consumer configuration
type KafkaConfig struct {
	Brokers            []string
	TopicProductUpdated string
	ConsumerGroup      string
	ReadTimeout        time.Duration
	MinBytes           int
	MaxBytes           int
}

// ElasticsearchConfig holds Elasticsearch connection configuration
type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
	IndexName string
	Timeout   time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level            string
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
}

// LoadConfig reads configuration from config.yaml and environment variables
// Environment variables take precedence over config file values
// Viper automatically maps environment variables (e.g., SERVER_PORT -> server.port)
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

	// Debug: Check viper values before unmarshal
	log.Printf("Viper values - ES index_name: %s, Kafka topic_product_updated: %s",
		viper.GetString("elasticsearch.index_name"),
		viper.GetString("kafka.topic_product_updated"),
	)

	// Unmarshal configuration into struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// If unmarshal didn't work, manually set values from viper
	if config.Elasticsearch.IndexName == "" {
		config.Elasticsearch.IndexName = viper.GetString("elasticsearch.index_name")
	}
	if config.Kafka.TopicProductUpdated == "" {
		config.Kafka.TopicProductUpdated = viper.GetString("kafka.topic_product_updated")
	}
	if config.Kafka.ConsumerGroup == "" {
		config.Kafka.ConsumerGroup = viper.GetString("kafka.consumer_group")
	}

	// Debug: Check if values were loaded
	log.Printf("After unmarshal - ES Index: %s, Kafka Topic: %s, ConsumerGroup: %s",
		config.Elasticsearch.IndexName,
		config.Kafka.TopicProductUpdated,
		config.Kafka.ConsumerGroup,
	)

	return config, nil
}

// setDefaults sets default values for configuration
// These are fallbacks if neither config file nor env vars are set
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8002)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic_product_updated", "product_updated")
	viper.SetDefault("kafka.consumer_group", "search-service")
	viper.SetDefault("kafka.read_timeout", "10s")
	viper.SetDefault("kafka.min_bytes", 1024)
	viper.SetDefault("kafka.max_bytes", 10485760) // 10MB

	// Elasticsearch defaults
	viper.SetDefault("elasticsearch.addresses", []string{"http://localhost:9200"})
	viper.SetDefault("elasticsearch.username", "")
	viper.SetDefault("elasticsearch.password", "")
	viper.SetDefault("elasticsearch.index_name", "products")
	viper.SetDefault("elasticsearch.timeout", "30s")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")
	viper.SetDefault("logging.output_paths", []string{"stdout"})
	viper.SetDefault("logging.error_output_paths", []string{"stderr"})
}

