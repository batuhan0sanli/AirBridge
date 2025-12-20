/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
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

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Decrypts a file received from a sender.",
	Long: `Starts an interactive session to receive a file.
It generates a temporary key pair for this session and displays the public key.
You can then provide the encrypted text block.`,
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

		p := tea.NewProgram(receive.InitialModel(initialPrivKeyPEM, initialPayload, inputPayloadPath, deletePayload))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
	receiveCmd.Flags().StringVarP(&privKeyPath, "privkey", "k", "", "Path to private key")
	receiveCmd.Flags().StringVarP(&inputPayloadPath, "input", "i", "", "Path to input payload file")
	receiveCmd.Flags().BoolVarP(&deletePayload, "delete", "d", false, "Delete payload file after successful decryption")
}
