package cmd

import (
	"fmt"
	"os"

	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/aes128-dev/aes128-cli/pkg/vpn"
	"github.com/spf13/cobra"
)

var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from the VPN",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() != 0 {
			fmt.Println("This command requires root privileges. Please run with sudo.")
			return
		}

		if err := vpn.Stop(); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Disconnected successfully.")
		}
		config.DeleteConnectionStatus()
	},
}