package auth

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/piotmni/better-auth-device-login/go-auth-device-cli/internal/config"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display authentication status",
	Long:  `Shows the current authentication state including user info and token expiry.`,
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		color.Red("Not logged in")
		return nil
	}

	// Check if we have tokens
	accessToken, err := config.GetAccessToken()
	if err != nil || accessToken == "" {
		color.Red("Not logged in")
		return nil
	}

	// Display status
	color.Green("Logged in")
	fmt.Printf("  Server: %s\n", cfg.Auth.Hostname)

	if cfg.Auth.UserEmail != "" {
		fmt.Printf("  User: %s\n", cfg.Auth.UserEmail)
	}

	if !cfg.Auth.ExpiresAt.IsZero() {
		if time.Now().After(cfg.Auth.ExpiresAt) {
			color.Yellow("  Token expired at: %s", cfg.Auth.ExpiresAt.Format(time.RFC3339))
		} else {
			fmt.Printf("  Token expires: %s\n", cfg.Auth.ExpiresAt.Format(time.RFC3339))
			fmt.Printf("  Time remaining: %s\n", time.Until(cfg.Auth.ExpiresAt).Round(time.Second))
		}
	}

	return nil
}
