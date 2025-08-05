package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const apiBaseURL = "https://client.aes128.com/api"

type AppSessionInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DnsSettingsResponse struct {
	RegularDNS string `json:"regular_dns"`
	AdblockDNS string `json:"adblock_dns"`
}

type LocationsResponse struct {
	UserUUID  string         `json:"user_uuid"`
	Locations []LocationInfo `json:"locations"`
}

type ApiResponse struct {
	AppSessionToken string           `json:"app_session_token"`
	SessionName     string           `json:"session_name"`
	Error           string           `json:"error"`
	Status          string           `json:"status"`
	Sessions        []AppSessionInfo `json:"sessions"`
}

type LocationInfo struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	IPAddress  string `json:"ip_address"`
	VlessPort  int    `json:"vless_port"`
	VmessPort  int    `json:"vmess_port"`
	TrojanPort int    `json:"trojan_port"`
}

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		token:      token,
	}
}

func (c *Client) sendRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("could not marshal payload: %w", err)
		}
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", apiBaseURL, endpoint), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-App-Session-Token", c.token)
	}

	return c.httpClient.Do(req)
}

func (c *Client) Login(username, password string) (*ApiResponse, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	resp, err := c.sendRequest("POST", "/app/login", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("could not parse server response: %w", err)
	}

	if resp.StatusCode == http.StatusConflict {
		apiResp.Error = "Maximum number of app sessions reached."
		return &apiResp, fmt.Errorf(apiResp.Error)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != "" {
			return nil, fmt.Errorf(apiResp.Error)
		}
		return nil, fmt.Errorf("invalid credentials or server error (status %d)", resp.StatusCode)
	}

	return &apiResp, nil
}

func (c *Client) Logout() error {
	resp, err := c.sendRequest("POST", "/app/logout", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status for logout: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetDnsSettings() (*DnsSettingsResponse, error) {
	resp, err := c.sendRequest("GET", "/app/dns_settings", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dnsResp DnsSettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&dnsResp); err != nil {
		return nil, fmt.Errorf("could not parse dns settings response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching dns settings (status %d)", resp.StatusCode)
	}

	return &dnsResp, nil
}

func (c *Client) GetLocations() (*LocationsResponse, error) {
	resp, err := c.sendRequest("GET", "/app/locations", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locResp LocationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&locResp); err != nil {
		return nil, fmt.Errorf("could not parse locations response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching locations (status %d)", resp.StatusCode)
	}

	return &locResp, nil
}

func (c *Client) GetSessions() ([]AppSessionInfo, error) {
	resp, err := c.sendRequest("GET", "/app/sessions", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sessions []AppSessionInfo
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return nil, fmt.Errorf("could not parse sessions response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching sessions (status %d)", resp.StatusCode)
	}

	return sessions, nil
}

func (c *Client) DeleteSession(sessionID int64) error {
	endpoint := "/app/sessions/delete/" + strconv.FormatInt(sessionID, 10)
	resp, err := c.sendRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete session (status %d)", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteSessionWithCredentials(username, password string, sessionIDToDelete int64) (*ApiResponse, error) {
	payload := map[string]interface{}{
		"username":             username,
		"password":             password,
		"session_id_to_delete": sessionIDToDelete,
	}
	resp, err := c.sendRequest("POST", "/app/delete-session", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("could not parse server response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || apiResp.Status != "success" {
		if apiResp.Error != "" {
			return nil, fmt.Errorf(apiResp.Error)
		}
		return nil, fmt.Errorf("failed to terminate session")
	}

	return &apiResp, nil
}