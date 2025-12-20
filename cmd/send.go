/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
	"AirBridge/internal/cli"
	"AirBridge/internal/tui/send"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var pubKeyPath string
var outputFilePath string
var headless bool

var sendCmd = &cobra.Command{
	Use:   "send [file]",
	Short: "Send a file securely.",
	Long: `Starts an interactive session to send a file.
It allows you to select a file, encrypt it with a recipient's public key,
and generates the encrypted payload.

You can optionally provide a file path as an argument to skip the file selection step.

Use --headless with -k and -o for headless mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		var initialFile string
		if len(args) > 0 {
			initialFile = args[0]
		}

		var initialPubKey string
		if pubKeyPath != "" {
			content, err := os.ReadFile(pubKeyPath)
			if err != nil {
				fmt.Printf("Error reading public key file: %v\n", err)
				os.Exit(1)
			}
			initialPubKey = string(content)
		}

		var appMode AppMode = ModeTUI
		if headless {
			appMode = ModeCLI
		}

		switch appMode {
		case ModeCLI:
			if initialFile == "" {
				fmt.Println("Error: File argument required in headless mode")
				os.Exit(1)
			}
			if initialPubKey == "" {
				fmt.Println("Error: Public key (-k) required in headless mode")
				os.Exit(1)
			}

			// Headless Execution
			if err := cli.RunSend(initialFile, initialPubKey, outputFilePath); err != nil {
				fmt.Printf("Error running headless send: %v\n", err)
				os.Exit(1)
			}

		case ModeTUI:
			p := tea.NewProgram(send.InitialModel(initialFile, initialPubKey, outputFilePath), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&pubKeyPath, "pubkey", "k", "", "Path to recipient's public key file (skips manual paste)")
	sendCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Path to save the payload file (default: payload.abp)")
	// Make the flag optional (NoOptDefVal) so -o works without an argument
	sendCmd.Flags().Lookup("output").NoOptDefVal = "payload.abp"
	sendCmd.Flags().BoolVarP(&headless, "headless", "H", false, "Run in headless mode (requires -k and file arg)")
}
