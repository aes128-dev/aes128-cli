package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/aes128-dev/aes128-cli/pkg/vpn"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "Show available VPN locations and their ping",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.ReadToken()
		if err != nil {
			fmt.Println("You are not logged in. Please run 'login' first.")
			return
		}

		fmt.Println("Fetching locations and checking ping...")

		client := api.NewClient(token)
		locationsResponse, err := client.GetLocations()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(locationsResponse.Locations) == 0 {
			fmt.Println("No locations available.")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Location", "Domain", "Ping (ms)"})
		
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetBorder(false) 
		table.SetHeaderLine(false)
		table.SetRowLine(false)
		table.SetColumnSeparator("   ") 
		table.SetCenterSeparator("")
		table.SetRowSeparator("")
		
		type pingResult struct {
			ID   int
			Data []string
			Ping int64
		}

		resultsChan := make(chan pingResult, len(locationsResponse.Locations))
		var wg sync.WaitGroup

		for i, loc := range locationsResponse.Locations {
			wg.Add(1)
			go func(location api.LocationInfo, id int) {
				defer wg.Done()
				var pingStr string
				var pingMs int64 = 9999

				ping, err := vpn.GetPing(location.IPAddress)
				if err != nil {
					pingStr = "N/A"
				} else {
					pingValue := ping.Milliseconds()
					if pingValue > 0 {
						pingMs = pingValue
						pingStr = strconv.FormatInt(pingMs, 10)
					} else if ping.Microseconds() > 0 {
						pingMs = 0
						pingStr = "0"
					} else {
						pingStr = "Timeout"
					}
				}
				resultsChan <- pingResult{ID: id, Data: []string{strconv.Itoa(id), location.Name, location.Domain, pingStr}, Ping: pingMs}
			}(loc, i+1)
		}

		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		var finalResults []pingResult
		for result := range resultsChan {
			finalResults = append(finalResults, result)
		}

		sort.Slice(finalResults, func(i, j int) bool {
			if finalResults[i].Ping != finalResults[j].Ping {
				return finalResults[i].Ping < finalResults[j].Ping
			}
			return finalResults[i].ID < finalResults[j].ID
		})
		
		table.Append([]string{})

		for _, res := range finalResults {
			table.Append(res.Data)
		}

		table.Render()
	},
}