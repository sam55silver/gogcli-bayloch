package bayloch

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the Bayloch dashboard connection settings.
type Config struct {
	URL    string // Dashboard base URL (e.g., "https://dashboard.bayloch.tech")
	APIKey string // Per-user API key (e.g., "bayloch_a1b2c3...")
}

// IsConfigured returns true if both GOG_BAYLOCH_URL and GOG_BAYLOCH_API_KEY are set.
func IsConfigured() bool {
	return strings.TrimSpace(os.Getenv("GOG_BAYLOCH_URL")) != "" &&
		strings.TrimSpace(os.Getenv("GOG_BAYLOCH_API_KEY")) != ""
}

// LoadConfig reads the Bayloch configuration from environment variables.
func LoadConfig() (*Config, error) {
	url := strings.TrimSpace(os.Getenv("GOG_BAYLOCH_URL"))
	if url == "" {
		return nil, fmt.Errorf("GOG_BAYLOCH_URL is not set")
	}
	url = strings.TrimRight(url, "/")

	apiKey := strings.TrimSpace(os.Getenv("GOG_BAYLOCH_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("GOG_BAYLOCH_API_KEY is not set")
	}

	return &Config{URL: url, APIKey: apiKey}, nil
}

// DashboardURL returns the dashboard URL for use in error messages.
func DashboardURL() string {
	url := strings.TrimSpace(os.Getenv("GOG_BAYLOCH_URL"))
	if url == "" {
		return "https://dashboard.bayloch.tech"
	}
	return strings.TrimRight(url, "/")
}
