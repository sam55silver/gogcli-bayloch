package bayloch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

// HTTPTokenSource implements oauth2.TokenSource by fetching tokens from the
// Bayloch dashboard HTTP API.
type HTTPTokenSource struct {
	cfg      *Config
	provider string
	email    string
	client   *http.Client
}

type tokenResponse struct {
	AccessToken   string   `json:"accessToken"`
	ProviderEmail string   `json:"providerEmail"`
	FolderID      string   `json:"folderId"`
	FolderName    string   `json:"folderName"`
	Scopes        []string `json:"scopes"`
	Error         string   `json:"error"`
}

// Token fetches a fresh access token from the Bayloch dashboard.
func (s *HTTPTokenSource) Token() (*oauth2.Token, error) {
	u := fmt.Sprintf("%s/api/tokens/%s", s.cfg.URL, s.provider)
	if s.email != "" {
		u += "?account=" + url.QueryEscape(s.email)
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("bayloch: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bayloch: dashboard unreachable at %s: %w", s.cfg.URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("bayloch: read response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// success
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("bayloch: invalid API key — find your key at %s", s.cfg.URL)
	case http.StatusNotFound:
		return nil, fmt.Errorf("bayloch: service %q not connected — connect it at %s", s.provider, s.cfg.URL)
	default:
		var errResp tokenResponse
		_ = json.Unmarshal(body, &errResp)
		msg := errResp.Error
		if msg == "" {
			msg = string(body)
		}
		return nil, fmt.Errorf("bayloch: unexpected status %d: %s", resp.StatusCode, msg)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("bayloch: parse response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("bayloch: empty access token in response")
	}

	return &oauth2.Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(55 * time.Minute),
	}, nil
}
