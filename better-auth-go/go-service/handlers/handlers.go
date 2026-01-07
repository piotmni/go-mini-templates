package handlers

import (
	"net/http"
	"time"

	"go-service/middleware"

	"github.com/labstack/echo/v4"
)

// HealthCheck returns the service health status
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Info returns public service information
func Info(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"service": "go-service",
		"version": "1.0.0",
		"message": "This is a public endpoint",
	})
}

// Profile returns the authenticated user's profile
func Profile(c echo.Context) error {
	ac := c.(*middleware.AuthContext)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    ac.UserID,
		"email": ac.Email,
	})
}

// SecureData returns protected data for authenticated users
func SecureData(c echo.Context) error {
	ac := c.(*middleware.AuthContext)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "This is protected data",
		"user_id": ac.UserID,
		"data": map[string]interface{}{
			"secret": "This data is only visible to authenticated users",
			"items":  []string{"item1", "item2", "item3"},
		},
	})
}
