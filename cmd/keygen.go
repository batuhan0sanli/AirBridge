/*
Copyright © 2025 Batuhan Sanli <batuhansanli@gmail.com>
*/
package cmd

import (
	"AirBridge/internal/crypto"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	outDir string
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new RSA key pair.",
	Long: `Generates a new RSA key pair (private.pem and public.pem)
in the specified directory (defaults to current directory).

These keys can be used for the send and receive commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating RSA key pair...")

		privateKey, publicKey, err := crypto.GenerateRSAKeyPair()
		if err != nil {
			fmt.Printf("Error generating keys: %v\n", err)
			os.Exit(1)
		}

		// Private Key
		privPEM, err := crypto.ExportRSAPrivateKeyAsPEM(privateKey)
		if err != nil {
			fmt.Printf("Error exporting private key: %v\n", err)
			os.Exit(1)
		}

		privateKeyPath := filepath.Join(outDir, "private.pem")
		if err := os.WriteFile(privateKeyPath, privPEM, 0600); err != nil {
			fmt.Printf("Error writing private key to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Private key saved to: %s\n", privateKeyPath)

		// Public Key
		pubPEM, err := crypto.ExportRSAPublicKeyAsPEM(publicKey)
		if err != nil {
			fmt.Printf("Error exporting public key: %v\n", err)
			os.Exit(1)
		}

		publicKeyPath := filepath.Join(outDir, "public.pem")
		if err := os.WriteFile(publicKeyPath, pubPEM, 0644); err != nil {
			fmt.Printf("Error writing public key to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Public key saved to: %s\n", publicKeyPath)
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	keygenCmd.Flags().StringVarP(&outDir, "output", "o", ".", "Directory to save the generated keys")
}
