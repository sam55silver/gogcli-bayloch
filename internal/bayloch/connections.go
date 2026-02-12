package bayloch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ConnectionInfo describes a connected service on the Bayloch dashboard.
type ConnectionInfo struct {
	Provider    string
	Email       string
	FolderName  string
	ConnectedAt time.Time
}

type connectionEntry struct {
	Email       *string   `json:"email"`
	FolderName  *string   `json:"folderName"`
	FolderID    *string   `json:"folderId"`
	ConnectedAt time.Time `json:"connectedAt"`
}

// ListConnections fetches the user's connected services from the dashboard.
func ListConnections() ([]ConnectionInfo, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/connections", cfg.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("bayloch: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bayloch: dashboard unreachable at %s: %w", cfg.URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("bayloch: read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("bayloch: invalid API key — find your key at %s", cfg.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bayloch: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var connectionsMap map[string]connectionEntry
	if err := json.Unmarshal(body, &connectionsMap); err != nil {
		return nil, fmt.Errorf("bayloch: parse response: %w", err)
	}

	connections := make([]ConnectionInfo, 0, len(connectionsMap))
	for provider, entry := range connectionsMap {
		info := ConnectionInfo{
			Provider:    provider,
			ConnectedAt: entry.ConnectedAt,
		}
		if entry.Email != nil {
			info.Email = *entry.Email
		}
		if entry.FolderName != nil {
			info.FolderName = *entry.FolderName
		}
		connections = append(connections, info)
	}

	return connections, nil
}
