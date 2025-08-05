package cmd

import (
	"fmt"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Manage active sessions (lists sessions by default)",
	Run: func(cmd *cobra.Command, args []string) {
		sessionsListCmd.Run(cmd, args)
	},
}

var sessionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active sessions",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.ReadToken()
		if err != nil {
			fmt.Println("You are not logged in.")
			return
		}
		client := api.NewClient(token)
		sessions, err := client.GetSessions()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(sessions) == 0 {
			fmt.Println("No active sessions found.")
			return
		}
		fmt.Println("Active sessions:")
		for _, s := range sessions {
			fmt.Printf("  - ID: %d, Name: %s\n", s.ID, s.Name)
		}
	},
}

var sessionsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an active session",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.ReadToken()
		if err != nil {
			fmt.Println("You are not logged in.")
			return
		}
		client := api.NewClient(token)
		sessions, err := client.GetSessions()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(sessions) == 0 {
			fmt.Println("No active sessions to delete.")
			return
		}

		var sessionNames []string
		for _, s := range sessions {
			sessionNames = append(sessionNames, fmt.Sprintf("ID: %d, Name: %s", s.ID, s.Name))
		}

		prompt := promptui.Select{
			Label: "Select session to delete",
			Items: sessionNames,
		}

		index, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		sessionToDelete := sessions[index]
		fmt.Printf("Deleting session: %s (ID: %d)...\n", sessionToDelete.Name, sessionToDelete.ID)

		err = client.DeleteSession(sessionToDelete.ID)
		if err != nil {
			fmt.Printf("Error deleting session: %v\n", err)
			return
		}
		fmt.Println("Session deleted successfully.")
	},
}