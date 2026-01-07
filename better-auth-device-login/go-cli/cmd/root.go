package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/piotmni/better-auth-device-login/go-auth-device-cli/cmd/auth"
)

var rootCmd = &cobra.Command{
	Use:   "go-auth-device-cli",
	Short: "CLI tool with device authorization authentication",
	Long:  `A CLI tool that authenticates with Better Auth using OAuth 2.0 Device Authorization Grant (RFC 8628).`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
}
