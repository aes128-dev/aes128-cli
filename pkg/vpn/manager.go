package vpn

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/aes128-dev/aes128-cli/pkg/api"
	"github.com/aes128-dev/aes128-cli/pkg/config"
)

const CorePath = "/usr/lib/aes128-cli/core"

const baseTemplate = `{
	"log": { "level": "info", "timestamp": true },
	"dns": { "servers": [], "strategy": "prefer_ipv4", "final": "dns-proxied" },
	"inbounds": [
		{ "type": "tun", "tag": "tun-in", "interface_name": "aes128tun", "stack": "system", "address": [ "172.19.0.1/30" ], "auto_route": true, "strict_route": true, "mtu": 1420 }
	],
	"outbounds": [
		{},
		{ "type": "direct", "tag": "direct" },
		{ "type": "block", "tag": "block" }
	],
	"route": {
		"auto_detect_interface": true,
		"rules": [
			{ "port": 53, "action": "hijack-dns" },
			{ "action": "sniff" },
			{ "ip_is_private": true, "outbound": "direct" }
		],
		"final": "proxy"
	}
}`

func FindFastestLocation(locations []api.LocationInfo) (string, error) {
	if len(locations) == 0 {
		return "", fmt.Errorf("location list is empty")
	}

	type pingResult struct {
		Domain string
		Rtt    time.Duration
	}

	resultsChan := make(chan pingResult, len(locations))
	var wg sync.WaitGroup

	for _, loc := range locations {
		wg.Add(1)
		go func(location api.LocationInfo) {
			defer wg.Done()
			ping, err := GetPing(location.IPAddress)
			if err == nil && ping > 0 {
				resultsChan <- pingResult{Domain: location.Domain, Rtt: ping}
			}
		}(loc)
	}

	wg.Wait()
	close(resultsChan)

	var bestDomain string
	minRtt := time.Hour

	for result := range resultsChan {
		if result.Rtt < minRtt {
			minRtt = result.Rtt
			bestDomain = result.Domain
		}
	}

	if bestDomain == "" {
		if len(locations) > 0 {
			return locations[0].Domain, nil
		}
		return "", fmt.Errorf("no responsive servers found")
	}

	return bestDomain, nil
}

func GenerateConfig(loc api.LocationInfo, settings *config.Settings, dns *api.DnsSettingsResponse, userUUID string) (string, error) {
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(baseTemplate), &configMap); err != nil {
		return "", fmt.Errorf("could not parse base config template: %w", err)
	}

	var outboundConfig map[string]interface{}
	switch settings.Protocol {
	case "vless":
		outboundConfig = map[string]interface{}{
			"type": "vless", "tag": "proxy", "server": loc.IPAddress, "server_port": loc.VlessPort, "uuid": userUUID,
			"tls": map[string]interface{}{"enabled": true, "server_name": loc.Domain},
		}
	case "vmess":
		outboundConfig = map[string]interface{}{
			"type": "vmess", "tag": "proxy", "server": loc.IPAddress, "server_port": loc.VmessPort, "uuid": userUUID, "security": "auto",
			"tls":       map[string]interface{}{"enabled": true, "server_name": loc.Domain},
			"transport": map[string]interface{}{"type": "ws", "path": "/vmess", "headers": map[string]string{"Host": loc.Domain}},
		}
	case "trojan":
		outboundConfig = map[string]interface{}{
			"type": "trojan", "tag": "proxy", "server": loc.IPAddress, "server_port": loc.TrojanPort, "password": userUUID,
			"tls": map[string]interface{}{"enabled": true, "server_name": loc.Domain},
		}
	default:
		return "", fmt.Errorf("unsupported protocol: %s", settings.Protocol)
	}
	configMap["outbounds"].([]interface{})[0] = outboundConfig

	dnsToUse := dns.RegularDNS
	if settings.AdBlock && dns.AdblockDNS != "" {
		dnsToUse = dns.AdblockDNS
	}
	dnsServers := []map[string]interface{}{
		{"tag": "dns-bootstrap", "address": "9.9.9.9", "detour": "direct"},
		{"tag": "dns-proxied", "address": dnsToUse, "address_resolver": "dns-bootstrap", "detour": "proxy"},
	}
	configMap["dns"].(map[string]interface{})["servers"] = dnsServers

	configBytes, err := json.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("could not marshal final config: %w", err)
	}
	return string(configBytes), nil
}

func Start(configContent string) error {
	pid, err := readPID()
	if err == nil && isProcessRunning(pid) {
		return fmt.Errorf("VPN is already running with PID %d", pid)
	}

	configPath, err := config.GetConfigFilePath(config.SingBoxConfigName)
	if err != nil {
		return fmt.Errorf("could not get config path: %w", err)
	}
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	cmd := exec.Command(CorePath, "run", "-c", configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start core process: %w", err)
	}
	fmt.Printf("VPN core process started with PID: %d\n", cmd.Process.Pid)
	return savePID(cmd.Process.Pid)
}

func Stop() error {
	pid, err := readPID()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("VPN is not running (no PID file found)")
		}
		return fmt.Errorf("could not read PID file: %w", err)
	}

	pidPath, _ := config.GetConfigFilePath(config.PIDFileName)

	if !isProcessRunning(pid) {
		os.Remove(pidPath)
		fmt.Printf("Cleaned up stale PID file for process %d.\n", pid)
		return fmt.Errorf("VPN is not running (process with PID %d not found)", pid)
	}

	fmt.Printf("Stopping VPN process with PID %d...\n", pid)
	cmd := exec.Command("kill", strconv.Itoa(pid))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}

	os.Remove(pidPath)
	return nil
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}

func savePID(pid int) error {
	path, err := config.GetConfigFilePath(config.PIDFileName)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
}

func readPID() (int, error) {
	path, err := config.GetConfigFilePath(config.PIDFileName)
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func GetConnectionDuration() (time.Duration, error) {
	pid, err := readPID()
	if err != nil {
		return 0, err
	}
	if !isProcessRunning(pid) {
		return 0, fmt.Errorf("process not running")
	}

	status, err := config.ReadConnectionStatus()
	if err != nil {
		return 0, err
	}
	return time.Since(status.StartTime), nil
}