package auth

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/piotmni/better-auth-device-login/go-auth-device-cli/internal/config"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the auth server",
	Long:  `Clears stored authentication tokens and removes auth configuration.`,
	RunE:  runLogout,
}

func runLogout(cmd *cobra.Command, args []string) error {
	// Clear tokens from keyring
	if err := config.ClearTokens(); err != nil {
		// Ignore errors if tokens don't exist
		fmt.Printf("Note: %v\n", err)
	}

	// Clear config
	if err := config.Clear(); err != nil {
		return fmt.Errorf("failed to clear config: %w", err)
	}

	color.Green("Logged out successfully")
	return nil
}
