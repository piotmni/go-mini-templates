package config

import (
	"github.com/zalando/go-keyring"
)

const (
	keyringService  = "go-auth-device-cli"
	accessTokenKey  = "access_token"
	refreshTokenKey = "refresh_token"
)

// SetAccessToken stores the access token in the system keyring
func SetAccessToken(token string) error {
	return keyring.Set(keyringService, accessTokenKey, token)
}

// GetAccessToken retrieves the access token from the system keyring
func GetAccessToken() (string, error) {
	return keyring.Get(keyringService, accessTokenKey)
}

// SetRefreshToken stores the refresh token in the system keyring
func SetRefreshToken(token string) error {
	return keyring.Set(keyringService, refreshTokenKey, token)
}

// GetRefreshToken retrieves the refresh token from the system keyring
func GetRefreshToken() (string, error) {
	return keyring.Get(keyringService, refreshTokenKey)
}

// ClearTokens removes all tokens from the system keyring
func ClearTokens() error {
	errAccess := keyring.Delete(keyringService, accessTokenKey)
	errRefresh := keyring.Delete(keyringService, refreshTokenKey)

	// Return first error if any
	if errAccess != nil && errAccess != keyring.ErrNotFound {
		return errAccess
	}
	if errRefresh != nil && errRefresh != keyring.ErrNotFound {
		return errRefresh
	}

	return nil
}
