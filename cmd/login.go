package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the service",
	Run:   runLogin,
}

func runLogin(cmd *cobra.Command, args []string) {
	if os.Geteuid() == 0 {
		fmt.Println("Error: do not run the 'login' command with sudo.")
		return
	}

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

	attemptLogin(username, password)
}

func attemptLogin(username, password string) {
	fmt.Println("Logging in...")
	client := api.NewClient("")
	result, err := client.Login(username, password)

	if err != nil {
		if result != nil && strings.Contains(err.Error(), "Maximum number of app sessions reached") {
			handleSessionLimit(username, password, result.Sessions)
		} else {
			fmt.Printf("Login error: %v\n", err)
		}
		return
	}

	saveTokenAndFetchData(result)
}

func handleSessionLimit(username, password string, sessions []api.AppSessionInfo) {
	fmt.Println("Device limit reached. Please choose a session to terminate.")
	if len(sessions) == 0 {
		fmt.Println("No sessions returned by server to delete.")
		return
	}

	var sessionNames []string
	for _, s := range sessions {
		sessionNames = append(sessionNames, s.Name)
	}

	prompt := promptui.Select{
		Label: "Select session to terminate",
		Items: sessionNames,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Selection cancelled: %v\n", err)
		return
	}

	sessionToDelete := sessions[index]
	fmt.Printf("Terminating session: %s...\n", sessionToDelete.Name)

	client := api.NewClient("")
	_, err = client.DeleteSessionWithCredentials(username, password, sessionToDelete.ID)
	if err != nil {
		fmt.Printf("Error terminating session: %v\n", err)
		return
	}

	fmt.Println("Session terminated. Retrying login...")
	attemptLogin(username, password)
}

func saveTokenAndFetchData(result *api.ApiResponse) {
	if result.AppSessionToken == "" {
		fmt.Println("Failed to get token from server.")
		return
	}
	if err := config.SaveToken(result.AppSessionToken); err != nil {
		fmt.Printf("Failed to save session token: %v\n", err)
		return
	}
	fmt.Println("Login successful. Fetching user data...")

	clientWithToken := api.NewClient(result.AppSessionToken)
	locationsResponse, err := clientWithToken.GetLocations()
	if err != nil {
		fmt.Printf("Could not fetch user data after login: %v\n", err)
		return
	}

	cache := &config.UserCache{
		UserUUID:  locationsResponse.UserUUID,
		Locations: locationsResponse.Locations,
	}

	if err := config.SaveUserCache(cache); err != nil {
		fmt.Printf("Could not save user data to cache: %v\n", err)
		return
	}

	fmt.Println("Session and user data saved.")
}