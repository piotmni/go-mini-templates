package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// Audience is a custom type that handles JWT audience claim which can be
// either a string or an array of strings per RFC 7519
type Audience []string

// UnmarshalJSON implements json.Unmarshaler to handle both string and []string
func (a *Audience) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single string first
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = Audience{single}
		return nil
	}

	// Try to unmarshal as an array of strings
	var multiple []string
	if err := json.Unmarshal(data, &multiple); err != nil {
		return err
	}
	*a = multiple
	return nil
}

// TokenClaims represents the claims extracted from the JWT token
type TokenClaims struct {
	ID        string   `json:"id,omitempty"`
	Email     string   `json:"email,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	Audience  Audience `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
}

// AuthContext extends echo.Context with authenticated user information
type AuthContext struct {
	echo.Context
	UserID string
	Email  string
}

// GetAuthContext extracts AuthContext from echo.Context
// Returns nil if the context is not an AuthContext (e.g., unauthenticated routes)
func GetAuthContext(c echo.Context) *AuthContext {
	if ac, ok := c.(*AuthContext); ok {
		return ac
	}
	return nil
}

var (
	jwkCache   *jwk.Cache
	cacheOnce  sync.Once
	cancelFunc context.CancelFunc
	initErr    error
)

// getAuthServiceURL returns the auth service base URL
func getAuthServiceURL() string {
	url := os.Getenv("AUTH_SERVICE_URL")
	if url == "" {
		url = "http://localhost:3000"
	}
	return url
}

// getJWKSURL returns the full JWKS endpoint URL
func getJWKSURL() string {
	return getAuthServiceURL() + "/api/auth/jwks"
}

// initCache initializes the JWK cache with auto-refresh capabilities
func initCache() (*jwk.Cache, error) {
	cacheOnce.Do(func() {
		// Create a background context for the cache lifecycle
		ctx, cancel := context.WithCancel(context.Background())
		cancelFunc = cancel

		// Create a new JWK cache with httprc client
		c, err := jwk.NewCache(ctx, httprc.NewClient())
		if err != nil {
			initErr = err
			return
		}

		// Register the JWKS URL with refresh intervals
		// MinInterval: minimum time between refreshes (even if server suggests shorter)
		// MaxInterval: maximum time between refreshes (even if server suggests longer)
		if err := c.Register(
			ctx,
			getJWKSURL(),
			jwk.WithMinInterval(1*time.Minute),
			jwk.WithMaxInterval(15*time.Minute),
		); err != nil {
			initErr = err
			return
		}

		jwkCache = c
	})
	return jwkCache, initErr
}

// ShutdownCache gracefully shuts down the JWK cache
// Should be called when the application is shutting down
func ShutdownCache() {
	if cancelFunc != nil {
		cancelFunc()
	}
}

// JWTAuth middleware validates JWT tokens using JWKS from the auth service
// Uses auto-refreshing cache to handle key rotation
// Wraps the context with AuthContext containing user information
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}
			tokenString := parts[1]

			// Get the cache (initialized once, thread-safe)
			cache, err := initCache()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize JWKS cache: "+err.Error())
			}

			// CachedSet returns a jwk.Set that transparently uses the cache
			cachedSet, err := cache.CachedSet(getJWKSURL())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get JWKS cache: "+err.Error())
			}

			// Parse and validate the token using the cached set
			// The cached set automatically uses the most recent keys
			token, err := jwt.Parse(
				[]byte(tokenString),
				jwt.WithKeySet(cachedSet),
				jwt.WithValidate(true),
			)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			// Marshal token to JSON and unmarshal into TokenClaims struct
			tokenBytes, err := json.Marshal(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal token: "+err.Error())
			}

			var claims TokenClaims
			if err := json.Unmarshal(tokenBytes, &claims); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal token claims: "+err.Error())
			}

			// Use subject as fallback for user ID if id claim is not set
			userID := claims.ID
			if userID == "" {
				userID = claims.Subject
			}

			// Wrap context with AuthContext containing user information
			ac := &AuthContext{
				Context: c,
				UserID:  userID,
				Email:   claims.Email,
			}

			return next(ac)
		}
	}
}
