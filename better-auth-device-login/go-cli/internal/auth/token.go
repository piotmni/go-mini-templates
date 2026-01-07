package auth

import "time"

// Token represents an OAuth token
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsExpired returns true if the token has expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid returns true if the token is valid and not expired
func (t *Token) IsValid() bool {
	return t.AccessToken != "" && !t.IsExpired()
}
