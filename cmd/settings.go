package cmd

import (
	"fmt"
	"strings"

	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/spf13/cobra"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "View or change settings",
	Long: `View current settings or change them.
Usage:
  aes128-cli settings get       (shows current settings)
  aes128-cli settings set <key> <value>  (sets a new value)

Examples:
  aes128-cli settings set protocol trojan
  aes128-cli settings set adblock on`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var settingsSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a new value for a setting",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])
		value := strings.ToLower(args[1])

		settings, err := config.LoadSettings()
		if err != nil {
			fmt.Printf("Error loading settings: %v\n", err)
			return
		}

		switch key {
		case "protocol":
			if value != "vless" && value != "vmess" && value != "trojan" {
				fmt.Println("Invalid protocol. Available: vless, vmess, trojan")
				return
			}
			settings.Protocol = value
			fmt.Printf("Protocol set to: %s\n", value)
		case "adblock":
			if value == "on" || value == "true" {
				settings.AdBlock = true
				fmt.Println("AdBlock enabled.")
			} else if value == "off" || value == "false" {
				settings.AdBlock = false
				fmt.Println("AdBlock disabled.")
			} else {
				fmt.Println("Invalid value for AdBlock. Use 'on' or 'off'.")
				return
			}
		default:
			fmt.Printf("Unknown setting key: %s\n", key)
			return
		}

		if err := config.SaveSettings(settings); err != nil {
			fmt.Printf("Error saving settings: %v\n", err)
		}
	},
}

var settingsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current settings",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.LoadSettings()
		if err != nil {
			fmt.Printf("Could not load settings: %v\n", err)
			return
		}
		fmt.Println("Current Settings:")
		fmt.Printf("  Protocol: %s\n", settings.Protocol)
		fmt.Printf("  AdBlock Enabled: %v\n", settings.AdBlock)
	},
}