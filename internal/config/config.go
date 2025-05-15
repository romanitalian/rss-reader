package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	RefreshInterval int      `json:"refresh_interval"` // Update interval in minutes
	AutoRefresh     bool     `json:"auto_refresh"`     // Automatic updates
	DefaultFeeds    []string `json:"default_feeds"`    // Default RSS feed URLs
	DataDir         string   `json:"data_dir"`         // Data storage directory
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Config{
		RefreshInterval: 30,
		AutoRefresh:     true,
		DefaultFeeds: []string{
			"https://blog.golang.org/feed.atom",
			"https://news.ycombinator.com/rss",
		},
		DataDir: filepath.Join(homeDir, ".rssreader"),
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	// If the file doesn't exist, return the default configuration
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		// Create the data directory if it doesn't exist
		if err := os.MkdirAll(config.DataDir, 0755); err != nil {
			return nil, err
		}
		return config, nil
	}

	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	// Create the data directory if it doesn't exist
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(configPath string, config *Config) error {
	// Create the configuration directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Serialize the configuration to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".rssreader", "config.json")
}
