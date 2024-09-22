// internal/server/config.go
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ServerConfig holds server-related configurations
type ServerConfig struct {
	Addr string `json:"addr"`
}

// DatabaseConfig holds database-related configurations
type DatabaseConfig struct {
	Driver     string `json:"driver"`
	DataSource string `json:"data_source"`
}

// Config holds the entire configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

// LoadConfig reads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}
	if cfg.Database.Driver != "sqlite" {
		return nil, fmt.Errorf("unsupported database driver: %s only sqlite is supported", cfg.Database.Driver)
	}
	//TODO: check if more to be added e.g valid filenames/paths.
	return &cfg, nil
}
