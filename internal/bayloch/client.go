package bayloch

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// TokenSourceForService returns an oauth2.TokenSource that fetches tokens from
// the Bayloch dashboard for the given gogcli service name.
// If email is non-empty, it selects a specific account via ?account=<email>.
func TokenSourceForService(service, email string) (oauth2.TokenSource, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	provider, err := ProviderForService(service)
	if err != nil {
		return nil, err
	}

	ts := &HTTPTokenSource{
		cfg:      cfg,
		provider: provider,
		email:    email,
		client:   &http.Client{Timeout: 10 * time.Second},
	}

	return oauth2.ReuseTokenSource(nil, ts), nil
}
