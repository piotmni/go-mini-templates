package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultClientID = "go-auth-device-cli"
	defaultScopes   = "openid profile email"
)

// Client is an HTTP client for auth endpoints
type Client struct {
	hostname   string
	httpClient *http.Client
	clientID   string
}

// NewClient creates a new auth client
func NewClient(hostname string) *Client {
	return &Client{
		hostname: hostname,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		clientID: defaultClientID,
	}
}

// DeviceCodeResponse represents the response from the device authorization endpoint
type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// TokenResponse represents the response from the token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenErrorResponse represents an error response from the token endpoint
type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// SessionResponse represents the session response
type SessionResponse struct {
	User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	} `json:"user"`
}

// RequestDeviceCode initiates the device authorization flow
func (c *Client) RequestDeviceCode(ctx context.Context) (*DeviceCodeResponse, error) {
	payload := map[string]string{
		"client_id": c.clientID,
		"scope":     defaultScopes,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.hostname+"/api/auth/device/code", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var deviceResp DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, err
	}

	// Set default interval if not provided
	if deviceResp.Interval == 0 {
		deviceResp.Interval = 5
	}

	return &deviceResp, nil
}

// PollForToken polls the token endpoint until authorization is complete or fails
func (c *Client) PollForToken(ctx context.Context, deviceCode string, interval int) (*TokenResponse, error) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			tokenResp, err := c.requestToken(ctx, deviceCode)
			if err == nil {
				return tokenResp, nil
			}

			// Check if it's a retryable error
			if tokenErr, ok := err.(*DeviceFlowError); ok {
				switch tokenErr.Code {
				case "authorization_pending":
					// Continue polling
					continue
				case "slow_down":
					// Increase interval
					ticker.Reset(time.Duration(interval+5) * time.Second)
					continue
				case "access_denied":
					return nil, fmt.Errorf("authorization denied by user")
				case "expired_token":
					return nil, fmt.Errorf("device code expired, please try again")
				}
			}

			return nil, err
		}
	}
}

// requestToken makes a single token request
func (c *Client) requestToken(ctx context.Context, deviceCode string) (*TokenResponse, error) {
	payload := map[string]string{
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		"device_code": deviceCode,
		"client_id":   c.clientID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.hostname+"/api/auth/device/token", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var tokenErr TokenErrorResponse
		if err := json.Unmarshal(respBody, &tokenErr); err == nil && tokenErr.Error != "" {
			return nil, &DeviceFlowError{
				Code:        tokenErr.Error,
				Description: tokenErr.ErrorDescription,
			}
		}
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetSession retrieves the current session information
func (c *Client) GetSession(ctx context.Context, accessToken string) (*SessionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.hostname+"/api/auth/session", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var session SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeviceFlowError represents an error in the device flow
type DeviceFlowError struct {
	Code        string
	Description string
}

func (e *DeviceFlowError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Description)
	}
	return e.Code
}
