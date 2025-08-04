package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiBaseURL = "https://client.aes128.com/api"

type AppSessionInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ApiResponse struct {
	AppSessionToken string           `json:"app_session_token"`
	SessionName     string           `json:"session_name"`
	Error           string           `json:"error"`
	Sessions        []AppSessionInfo `json:"sessions"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Login(username, password string) (*ApiResponse, error) {
	payload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/app/login", apiBaseURL), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not connect to authentication service: %w", err)
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("could not parse server response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != "" {
			return nil, fmt.Errorf(apiResp.Error)
		}
		if resp.StatusCode == http.StatusConflict {
			return &apiResp, fmt.Errorf("device limit reached")
		}
		return nil, fmt.Errorf("invalid credentials or server error (status %d)", resp.StatusCode)
	}

	return &apiResp, nil
}