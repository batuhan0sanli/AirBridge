/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type AppMode int

const (
	ModeTUI AppMode = iota
	ModeCLI
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
	Version: "v0.1.3",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
