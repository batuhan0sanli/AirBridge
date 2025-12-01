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
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Decrypts a file received from a sender.",
	Long: `Starts an interactive session to receive a file.
It generates a temporary key pair for this session and displays the public key.
You can then provide the encrypted text block.`, // or the path to a .abf file.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(receive.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
