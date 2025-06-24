package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Database configuration
	Database DatabaseConfig

	// Server configuration
	Server ServerConfig

	// Kafka configuration
	Kafka KafkaConfig

	// Environment
	Environment string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	URL      string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers      string
	Topic        string
	ZookeeperURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Load database configuration
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "piush"),
		Password: getEnv("DB_PASSWORD", "root_access"),
		Name:     getEnv("DB_NAME", "everato"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		URL:      getEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?sslmode=disable"),
	}

	// Load server configuration
	serverPort, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	config.Server = ServerConfig{
		Port: serverPort,
	}

	// Load Kafka configuration
	config.Kafka = KafkaConfig{
		Brokers:      getEnv("KAFKA_BROKERS", "localhost:9092"),
		Topic:        getEnv("KAFKA_TOPIC", "everato-events"),
		ZookeeperURL: getEnv("ZOOKEEPER_URL", "localhost:2181"),
	}

	// Load environment
	config.Environment = getEnv("ENV", "development")

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// DatabaseURL returns the complete database URL
func (c *Config) DatabaseURL() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
