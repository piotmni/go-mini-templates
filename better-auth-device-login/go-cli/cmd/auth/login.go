package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	internalAuth "github.com/piotmni/better-auth-device-login/go-auth-device-cli/internal/auth"
	"github.com/piotmni/better-auth-device-login/go-auth-device-cli/internal/browser"
	"github.com/piotmni/better-auth-device-login/go-auth-device-cli/internal/config"
)

var (
	loginHostname  string
	loginNoBrowser bool
	loginTimeout   time.Duration
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the auth server",
	Long:  `Initiates OAuth 2.0 Device Authorization Grant flow to authenticate with the auth server.`,
	RunE:  runLogin,
}

func init() {
	loginCmd.Flags().StringVar(&loginHostname, "hostname", "", "Auth server URL")
	loginCmd.Flags().BoolVar(&loginNoBrowser, "no-browser", false, "Don't open browser automatically")
	loginCmd.Flags().DurationVar(&loginTimeout, "timeout", 5*time.Minute, "Polling timeout")
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Get hostname from flag, config, or prompt
	hostname := loginHostname
	if hostname == "" {
		cfg, _ := config.Load()
		if cfg != nil && cfg.Auth.Hostname != "" {
			hostname = cfg.Auth.Hostname
		}
	}
	if hostname == "" {
		hostname = promptForHostname()
	}

	// Normalize hostname
	hostname = strings.TrimSuffix(hostname, "/")

	fmt.Printf("Authenticating with %s\n", hostname)

	// Create auth client
	client := internalAuth.NewClient(hostname)

	// Request device code
	deviceResp, err := client.RequestDeviceCode(context.Background())
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}

	// Display instructions
	fmt.Println()
	color.Cyan("Please visit: %s", deviceResp.VerificationURI)
	fmt.Printf("And enter code: ")
	color.Yellow(deviceResp.UserCode)
	fmt.Println()

	// Open browser if not disabled
	if !loginNoBrowser {
		uri := deviceResp.VerificationURIComplete
		if uri == "" {
			uri = deviceResp.VerificationURI
		}
		if err := browser.Open(uri); err != nil {
			fmt.Printf("Could not open browser: %v\n", err)
			fmt.Println("Please open the URL manually.")
		}
	}

	// Start polling with spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Waiting for authorization..."
	s.Start()

	ctx, cancel := context.WithTimeout(context.Background(), loginTimeout)
	defer cancel()

	tokenResp, err := client.PollForToken(ctx, deviceResp.DeviceCode, deviceResp.Interval)
	s.Stop()

	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Store tokens in keyring
	if err := config.SetAccessToken(tokenResp.AccessToken); err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}
	if tokenResp.RefreshToken != "" {
		if err := config.SetRefreshToken(tokenResp.RefreshToken); err != nil {
			return fmt.Errorf("failed to store refresh token: %w", err)
		}
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Get user info if possible
	userEmail := ""
	session, err := client.GetSession(context.Background(), tokenResp.AccessToken)
	if err == nil && session != nil && session.User.Email != "" {
		userEmail = session.User.Email
	}

	// Save config
	cfg := &config.Config{
		Auth: config.AuthConfig{
			Hostname:  hostname,
			UserEmail: userEmail,
			ExpiresAt: expiresAt,
		},
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	color.Green("Authentication successful!")
	if userEmail != "" {
		fmt.Printf("Logged in as: %s\n", userEmail)
	}

	return nil
}

func promptForHostname() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter auth server URL (e.g., http://localhost:3000): ")
	hostname, _ := reader.ReadString('\n')
	return strings.TrimSpace(hostname)
}
