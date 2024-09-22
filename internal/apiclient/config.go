// internal/apiclient/config.go
package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ClientConfig holds the client configuration values
type ClientConfig struct {
	Username   string `json:"username"`
	BackendURL string `json:"backend_url"`
}

// LoadConfig reads configuration from a JSON file
func LoadConfig(path string) (*ClientConfig, error) {
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

	var cfg ClientConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if strings.TrimSpace(cfg.Username) == "" || !isValidURL(cfg.BackendURL) {
		return nil, fmt.Errorf("username and backend_url must be set in the config file")
	}

	return &cfg, nil
}

// DefaultConfigPath returns the default config file path
func DefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "team-timetracker.json")
}

// isValidURL checks if a given URL string is valid and matches the expected
// requirements (e.g. scheme, host).
func isValidURL(theURL string) bool {
	if strings.TrimSpace(theURL) == "" {
		return false
	}

	u, err := url.Parse(theURL)
	if err != nil {
		return false
	}

	// Additional checks for specific requirements (optional)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
