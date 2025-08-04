package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aes128-cli",
	Short: "aes128-cli is a command-line VPN client",
	Long:  `A fast and reliable CLI for the AES128 VPN service.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}