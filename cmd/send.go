/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
	"AirBridge/internal/tui/send"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [file]",
	Short: "Send a file securely.",
	Long: `Starts an interactive session to send a file.
It allows you to select a file, encrypt it with a recipient's public key,
and generates the encrypted payload.

You can optionally provide a file path as an argument to skip the file selection step.`,
	Run: func(cmd *cobra.Command, args []string) {
		var initialFile string
		if len(args) > 0 {
			initialFile = args[0]
		}
		p := tea.NewProgram(send.InitialModel(initialFile), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
