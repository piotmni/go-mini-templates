package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/piotmni/better-auth-device-login/go-auth-device-cli/internal/config"
)

// Client is an authenticated HTTP client for API requests
type Client struct {
	hostname   string
	httpClient *http.Client
}

// NewClient creates a new authenticated API client
func NewClient() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("not logged in: %w", err)
	}

	return &Client{
		hostname: cfg.Auth.Hostname,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Do performs an authenticated HTTP request
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Get access token from keyring
	accessToken, err := config.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return c.httpClient.Do(req)
}

// NewRequest creates a new HTTP request with the base URL
func (c *Client) NewRequest(ctx context.Context, method, path string) (*http.Request, error) {
	url := c.hostname + path
	return http.NewRequestWithContext(ctx, method, url, nil)
}

// Hostname returns the configured auth server hostname
func (c *Client) Hostname() string {
	return c.hostname
}
