/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
	"AirBridge/internal/cli"
	"AirBridge/internal/tui/receive"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// receiveCmd represents the receive command
var privKeyPath string
var inputPayloadPath string
var deletePayload bool
var headlessReceive bool

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Decrypts a file received from a sender.",
	Long: `Starts an interactive session to receive a file.
It generates a temporary key pair for this session and displays the public key.
You can then provide the encrypted text block.

Use --headless with -k and -i for headless mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		var initialPrivKeyPEM []byte
		if privKeyPath != "" {
			var err error
			initialPrivKeyPEM, err = os.ReadFile(privKeyPath)
			if err != nil {
				fmt.Printf("Error reading private key file: %v\n", err)
				os.Exit(1)
			}
		}

		var initialPayload string
		if inputPayloadPath != "" {
			content, err := os.ReadFile(inputPayloadPath)
			if err != nil {
				fmt.Printf("Error reading payload file: %v\n", err)
				os.Exit(1)
			}
			initialPayload = string(content)
		}

		var appMode AppMode = ModeTUI
		if headlessReceive {
			appMode = ModeCLI
		}

		switch appMode {
		case ModeCLI:
			if len(initialPrivKeyPEM) == 0 {
				fmt.Println("Error: Private key (-k) required in headless mode")
				os.Exit(1)
			}
			if initialPayload == "" {
				fmt.Println("Error: Input payload (-i) required in headless mode")
				os.Exit(1)
			}

			// Headless Execution
			if err := cli.RunReceive(initialPayload, initialPrivKeyPEM, inputPayloadPath, deletePayload); err != nil {
				fmt.Printf("Error running headless receive: %v\n", err)
				os.Exit(1)
			}

		case ModeTUI:
			p := tea.NewProgram(receive.InitialModel(initialPrivKeyPEM, initialPayload, inputPayloadPath, deletePayload))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
	receiveCmd.Flags().StringVarP(&privKeyPath, "privkey", "k", "", "Path to private key")
	receiveCmd.Flags().StringVarP(&inputPayloadPath, "input", "i", "", "Path to input payload file")
	receiveCmd.Flags().BoolVarP(&deletePayload, "delete", "d", false, "Delete payload file after successful decryption")
	receiveCmd.Flags().BoolVarP(&headlessReceive, "headless", "H", false, "Run in headless mode (requires -k and -i)")
}
