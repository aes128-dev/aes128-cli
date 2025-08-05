package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/aes128-dev/aes128-cli/pkg/vpn"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current connection status",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() != 0 {
			fmt.Println("This command requires root privileges to check the service status. Please run with sudo.")
			return
		}

		status, err := config.ReadConnectionStatus()
		if err != nil {
			fmt.Println("Status: Disconnected")
			return
		}

		duration, err := vpn.GetConnectionDuration()
		if err != nil {
			fmt.Println("Status: Disconnected (VPN process not found or has been terminated)")
			config.DeleteConnectionStatus()
			return
		}

		fmt.Println("Status: Connected")
		fmt.Printf("Location: %s\n", status.LocationName)
		fmt.Printf("Uptime: %s\n", formatDuration(duration))
	},
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
