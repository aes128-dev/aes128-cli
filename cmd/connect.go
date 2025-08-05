package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
	"github.com/aes128-dev/aes128-cli/pkg/vpn"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [location_id_or_domain]",
	Short: "Connect to a VPN location (defaults to the fastest)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() != 0 {
			fmt.Println("This command requires root privileges. Please run with sudo.")
			return
		}

		if err := downloadAndInstallCore(); err != nil {
			fmt.Printf("Error during core setup: %v\n", err)
			return
		}

		token, err := config.ReadToken()
		if err != nil {
			fmt.Println("You are not logged in. Please run 'login' first.")
			return
		}
		client := api.NewClient(token)

		fmt.Println("Reading user cache...")
		cache, err := config.ReadUserCache()
		if err != nil {
			fmt.Printf("Could not read user cache. Please try to log in again. Error: %v\n", err)
			return
		}

		var target string
		if len(args) == 0 {
			fmt.Println("No location specified, finding the fastest server...")
			fastestDomain, err := vpn.FindFastestLocation(cache.Locations)
			if err != nil {
				fmt.Printf("Could not find the fastest server: %v\n", err)
				return
			}
			fmt.Printf("Fastest server found: %s\n", fastestDomain)
			target = fastestDomain
		} else {
			target = args[0]
		}

		var targetLocation *api.LocationInfo
		for i, loc := range cache.Locations {
			if strconv.Itoa(i+1) == target || strings.EqualFold(loc.Domain, target) {
				l := loc
				targetLocation = &l
				break
			}
		}

		if targetLocation == nil {
			fmt.Printf("Could not find location matching '%s'. Run 'locations' to see the list.\n", target)
			return
		}

		fmt.Printf("Selected location: %s\n", targetLocation.Name)

		fmt.Println("Fetching settings...")
		settings, err := config.LoadSettings()
		if err != nil {
			fmt.Printf("Error loading settings: %v\n", err)
			return
		}
		dns, err := client.GetDnsSettings()
		if err != nil {
			fmt.Printf("Error fetching DNS settings: %v\n", err)
			return
		}

		userUUID := cache.UserUUID
		if userUUID == "" {
			fmt.Println("User UUID not found in cache. Please log in again.")
			return
		}

		fmt.Println("Generating configuration...")
		configString, err := vpn.GenerateConfig(*targetLocation, settings, dns, userUUID)
		if err != nil {
			fmt.Printf("Error generating config: %v\n", err)
			return
		}

		fmt.Printf("Connecting to %s via %s protocol...\n", targetLocation.Name, settings.Protocol)
		if err := vpn.Start(configString); err != nil {
			fmt.Printf("Connection failed: %v\n", err)
			return
		}

		status := &config.ConnectionStatus{
			LocationName: targetLocation.Name,
			StartTime:    time.Now(),
		}
		if err := config.SaveConnectionStatus(status); err != nil {
			fmt.Printf("Warning: could not save connection status: %v\n", err)
		}

		fmt.Println("\nConnection successful!")
	},
}

func downloadAndInstallCore() error {
	if _, err := os.Stat(vpn.CorePath); err == nil {
		return nil
	}

	fmt.Println("VPN core not found. Attempting to download...")
	arch := runtime.GOARCH
	singboxVersion := "1.11.15"

	var downloadURL string
	switch arch {
	case "amd64":
		downloadURL = fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/sing-box-%s-linux-amd64.tar.gz", singboxVersion, singboxVersion)
	case "arm64":
		downloadURL = fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/sing-box-%s-linux-arm64.tar.gz", singboxVersion, singboxVersion)
	default:
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	fmt.Printf("Downloading for architecture %s...\n", arch)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)
	
	installDir := "/usr/lib/aes128-cli"
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}
	
	found := false
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read from tar archive: %w", err)
		}
		
		if strings.HasSuffix(header.Name, "sing-box") && header.Typeflag == tar.TypeReg {
			outFile, err := os.OpenFile(vpn.CorePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create core file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to extract core file: %w", err)
			}
			outFile.Close()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("sing-box binary not found in archive")
	}

	fmt.Println("VPN core installed successfully.")
	return nil
}