package cmd

import (
	"fmt"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the service and save session token",
	Run: func(cmd *cobra.Command, args []string) {
		promptUser := promptui.Prompt{
			Label: "Identifier",
		}
		username, err := promptUser.Run()
		if err != nil {
			fmt.Printf("Input cancelled: %v\n", err)
			return
		}

		promptPass := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}
		password, err := promptPass.Run()
		if err != nil {
			fmt.Printf("Input cancelled: %v\n", err)
			return
		}

		fmt.Println("Logging in...")

		client := api.NewClient()
		result, err := client.Login(username, password)
		if err != nil {
			if result != nil && len(result.Sessions) > 0 {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Please remove an active session on the aes128.com website.")
				return
			}
			fmt.Printf("Login error: %v\n", err)
			return
		}

		if result.AppSessionToken == "" {
			fmt.Println("Failed to get token from server.")
			return
		}

		if err := config.SaveToken(result.AppSessionToken); err != nil {
			fmt.Printf("Failed to save session token: %v\n", err)
			return
		}

		fmt.Println("Login successful. Session saved.")
	},
}