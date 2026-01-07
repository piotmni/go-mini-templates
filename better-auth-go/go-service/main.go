package main

import (
	"go-service/handlers"
	jwtmiddleware "go-service/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes
	public := e.Group("/api/public")
	public.GET("/health", handlers.HealthCheck)
	public.GET("/info", handlers.Info)

	// Protected routes (require JWT)
	protected := e.Group("/api/protected")
	protected.Use(jwtmiddleware.JWTAuth())
	protected.GET("/profile", handlers.Profile)
	protected.GET("/data", handlers.SecureData)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
