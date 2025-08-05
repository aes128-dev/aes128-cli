package cmd

import (
	"fmt"
	"os"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the service",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() == 0 {
			fmt.Println("Error: do not run the 'logout' command with sudo.")
			return
		}

		token, err := config.ReadToken()
		if err != nil {
			fmt.Println("You are already logged out.")
			return
		}

		pidPath, err := config.GetConfigFilePath(config.PIDFileName)
		if err == nil {
			if _, err := os.Stat(pidPath); err == nil {
				fmt.Println("Error: An active VPN connection was found.")
				fmt.Println("You must disconnect from the VPN before logging out.")
				fmt.Println("Please run 'sudo aes128-cli disconnect' first.")
				return
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
