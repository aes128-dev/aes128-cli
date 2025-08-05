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

func init() {
	cobra.EnableCommandSorting = false

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

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(locationsCmd)
	rootCmd.AddCommand(settingsCmd)

	connectCmd.GroupID = "core"
	disconnectCmd.GroupID = "core"
	statusCmd.GroupID = "core"

	loginCmd.GroupID = "management"
	logoutCmd.GroupID = "management"
	locationsCmd.GroupID = "management"
	settingsCmd.GroupID = "management"
}