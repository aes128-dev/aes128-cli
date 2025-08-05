package cmd

import (
	"fmt"
	"os"

	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Show current account information",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() == 0 {
			fmt.Println("Error: do not run this command with sudo.")
			return
		}

		cache, err := config.ReadUserCache()
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("You are not logged in. Please run 'login' first.")
			} else {
				fmt.Printf("Could not read user cache. Please try to log in again. Error: %v\n", err)
			}
			return
		}

		if cache.Username == "" || cache.SessionName == "" {
			fmt.Println("Account information is incomplete. Please try logging in again.")
			return
		}

		fmt.Println("Current Account Info:")
		fmt.Printf("  Username: %s\n", cache.Username)
		fmt.Printf("  Session Name: %s\n", cache.SessionName)
	},
}
