/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "airbridge",
	Short: "AirBridge is a secure, TUI-based file transfer tool.",
	Long: `AirBridge is a secure file transfer tool that uses RSA encryption to safely share files.
It features a TUI (Terminal User Interface) for easy interaction, allowing you to:
- Generate and share public keys.
- Encrypt files with a recipient's public key.
- Decrypt received payloads using your private key.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
