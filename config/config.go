package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name           string      `yaml:"name"`
	Version        string      `yaml:"version"`
	ApiPrefix      string      `yaml:"api_prefix"`
	Logging        bool        `yaml:"logging"`
	RequestTimeout string      `yaml:"req_timeout"`
	Server         Server      `yaml:"server"`
	DataBase       DataBase    `yaml:"database"`
	SupserUsers    []SuperUser `yaml:"super_users"`
}

type SuperUser struct {
	Name        string   `yaml:"name"`
	Password    string   `yaml:"password"`
	Email       string   `yaml:"email"`
	UserName    string   `yaml:"username"`
	Permissions []string `yaml:"permissions"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	TLS  TLS    `yaml:"tls"`
}

type TLS struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// Loads and returns the new config mapped into a
// struct of type Config
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

// PrettyPrint prints the given value in a pretty JSON format.
func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}
