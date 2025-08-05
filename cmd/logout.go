package cmd

import (
	"fmt"
	"os"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/aes128-dev/aes128-cli/pkg/vpn"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the service",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() == 0 {
			fmt.Println("Error: do not run the 'logout' command with sudo.")
			fmt.Println("Logout command works with your user files and does not require root privileges.")
			return
		}

		token, err := config.ReadToken()
		if err != nil {
			fmt.Println("You are already logged out.")
			return
		}

		pidPath, err := config.GetConfigFilePath(config.PIDFileName)
		if err != nil {
			fmt.Println("Warning: could not get PID file path.")
		}

		if _, err := os.Stat(pidPath); err == nil {
			fmt.Println("Active VPN connection found. Disconnecting first...")
			if err := vpn.Stop(); err != nil {
				fmt.Printf("Could not disconnect cleanly, but proceeding with logout. Error: %v\n", err)
			} else {
				fmt.Println("Disconnected successfully.")
			}
		}

		fmt.Println("Logging out from server...")
		client := api.NewClient(token)
		if err := client.Logout(); err != nil {
			fmt.Printf("Warning: could not log out from server. You may need to terminate this session manually via the website. Error: %v\n", err)
		}

		config.ClearSessionData()
		fmt.Println("Local session data cleared.")
		fmt.Println("Logout successful.")
	},
}