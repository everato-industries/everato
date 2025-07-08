// Package config provides functionality for loading and managing application configuration.
// It supports loading configuration from YAML files and provides structured access to
// configuration values for server settings, database connections, and user permissions.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the main application configuration structure.
// It contains all settings loaded from the configuration file.
type Config struct {
	Name           string      `yaml:"name"`        // Application name
	Version        string      `yaml:"version"`     // Application version
	ApiPrefix      string      `yaml:"api_prefix"`  // Prefix for all API routes (e.g., /api/v1)
	Logging        bool        `yaml:"logging"`     // Whether to enable logging
	RequestTimeout string      `yaml:"req_timeout"` // Maximum request timeout duration (e.g., "15s")
	Server         Server      `yaml:"server"`      // Server-specific configuration
	DataBase       DataBase    `yaml:"database"`    // Database connection configuration
	SupserUsers    []SuperUser `yaml:"super_users"` // List of predefined super users
}

// SuperUser represents a predefined administrative user with elevated permissions.
// These users are typically created during initial application setup.
type SuperUser struct {
	Name        string   `yaml:"name"`        // Full name of the super user
	Password    string   `yaml:"password"`    // Initial password (should be hashed in production)
	Email       string   `yaml:"email"`       // Email address for the super user
	UserName    string   `yaml:"username"`    // Username for login
	Permissions []string `yaml:"permissions"` // List of permission identifiers granted to this user
}

// DataBase contains configuration for database connection.
// This structure holds all parameters needed to establish a connection to the database.
type DataBase struct {
	Host     string `yaml:"host"`     // Database server hostname or IP
	Port     int    `yaml:"port"`     // Database server port
	User     string `yaml:"user"`     // Database username
	Password string `yaml:"password"` // Database password
	Name     string `yaml:"name"`     // Database name to connect to
}

// Server contains HTTP server configuration parameters.
// This includes listening address, port, and TLS settings.
type Server struct {
	Host string `yaml:"host"` // Server host address (e.g., "0.0.0.0" to listen on all interfaces)
	Port int    `yaml:"port"` // Server port number
	TLS  TLS    `yaml:"tls"`  // TLS/HTTPS configuration
}

// TLS contains configuration for HTTPS/TLS server settings.
// When enabled, the server will use these settings for secure connections.
type TLS struct {
	Enabled  bool   `yaml:"enabled"`   // Whether TLS is enabled
	CertFile string `yaml:"cert_file"` // Path to TLS certificate file
	KeyFile  string `yaml:"key_file"`  // Path to TLS private key file
}

// NewConfig loads configuration from a YAML file at the specified path.
// It reads the file, parses its contents, and returns a populated Config struct.
//
// Parameters:
//   - path: The path to the configuration file. If empty, defaults to "config.yaml".
//
// Returns:
//   - A pointer to the loaded Config structure
//   - An error if the file cannot be found, read, or parsed
func NewConfig(path string) (*Config, error) {
	config := &Config{}

	if path == "" {
		path = "config.yaml"
	}

	// Check if the file at the path exists or not
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at path: %s", path)
	}

	// Read the file data into a string
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the YAML data into the Config struct
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %v", err)
	}

	return config, nil
}

// PrettyPrint outputs the provided value as formatted JSON to standard output.
// This is useful for debugging and displaying configuration values.
//
// Parameters:
//   - v: Any value that can be marshalled to JSON
func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}
