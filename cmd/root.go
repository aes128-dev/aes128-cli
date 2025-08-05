package cmd

import (
	"fmt"
	"os"

	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aes128-cli",
	Short: "aes128-cli is a command-line VPN client",
	Long:  `A fast and reliable CLI for the AES128 VPN service.`,
}

func Execute() {
	setupCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupCommands() {
	cobra.EnableCommandSorting = false

	_, err := config.ReadToken()
	if err != nil {
		rootCmd.AddCommand(loginCmd)
	} else {
		rootCmd.AddGroup(&cobra.Group{
			ID:    "core",
			Title: "Core Commands:",
		})
		rootCmd.AddGroup(&cobra.Group{
			ID:    "management",
			Title: "Management Commands:",
		})

		rootCmd.AddCommand(connectCmd)
		rootCmd.AddCommand(disconnectCmd)
		rootCmd.AddCommand(statusCmd)
		rootCmd.AddCommand(locationsCmd)
		rootCmd.AddCommand(settingsCmd)
		rootCmd.AddCommand(accountCmd)
		rootCmd.AddCommand(logoutCmd)

		connectCmd.GroupID = "core"
		disconnectCmd.GroupID = "core"
		statusCmd.GroupID = "core"

		locationsCmd.GroupID = "management"
		settingsCmd.GroupID = "management"
		accountCmd.GroupID = "management"
		logoutCmd.GroupID = "management"
	}
}
