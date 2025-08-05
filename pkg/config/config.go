package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aes128-dev/aes128-cli/pkg/api"
)

const (
	ConfigDirName     = "aes128-cli"
	TokenFileName     = "session.token"
	SettingsFileName  = "settings.json"
	StatusFileName    = "status.json"
	PIDFileName       = "core.pid"
	SingBoxConfigName = "singbox.json"
	CacheFileName     = "cache.json"
)

type Settings struct {
	Protocol string `json:"protocol"`
	AdBlock  bool   `json:"adblock"`
}

type ConnectionStatus struct {
	LocationName string    `json:"locationName"`
	StartTime    time.Time `json:"startTime"`
}

type UserCache struct {
	UserUUID  string            `json:"user_uuid"`
	Locations []api.LocationInfo `json:"locations"`
}

func getHomeDir() (string, error) {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		currentUser, err := user.Lookup(sudoUser)
		if err != nil {
			return "", err
		}
		return currentUser.HomeDir, nil
	}
	return os.UserHomeDir()
}

func GetConfigFilePath(fileName string) (string, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config")
	appConfigDir := filepath.Join(configDir, ConfigDirName)

	if err := os.MkdirAll(appConfigDir, 0750); err != nil {
		return "", err
	}

	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		currentUser, err := user.Lookup(sudoUser)
		if err == nil {
			uid, _ := strconv.Atoi(currentUser.Uid)
			gid, _ := strconv.Atoi(currentUser.Gid)
			os.Chown(appConfigDir, uid, gid)
		}
	}

	return filepath.Join(appConfigDir, fileName), nil
}

func removeFile(fileName string) {
	path, err := GetConfigFilePath(fileName)
	if err == nil {
		os.Remove(path)
	}
}

func ClearSessionData() {
	removeFile(TokenFileName)
	removeFile(CacheFileName)
	removeFile(StatusFileName)
	removeFile(PIDFileName)
}

func SaveToken(token string) error {
	path, err := GetConfigFilePath(TokenFileName)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(token), 0600)
}

func ReadToken() (string, error) {
	path, err := GetConfigFilePath(TokenFileName)
	if err != nil {
		return "", err
	}
	token, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func LoadSettings() (*Settings, error) {
	path, err := GetConfigFilePath(SettingsFileName)
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			defaultSettings := &Settings{Protocol: "vless", AdBlock: false}
			if err := SaveSettings(defaultSettings); err != nil {
				return nil, err
			}
			return defaultSettings, nil
		}
		return nil, err
	}
	var settings Settings
	if err := json.Unmarshal(file, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

func SaveSettings(settings *Settings) error {
	path, err := GetConfigFilePath(SettingsFileName)
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0600)
}

func SaveConnectionStatus(status *ConnectionStatus) error {
	path, err := GetConfigFilePath(StatusFileName)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func ReadConnectionStatus() (*ConnectionStatus, error) {
	path, err := GetConfigFilePath(StatusFileName)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var status ConnectionStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func DeleteConnectionStatus() error {
	path, err := GetConfigFilePath(StatusFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(path)
}

func SaveUserCache(cache *UserCache) error {
	path, err := GetConfigFilePath(CacheFileName)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func ReadUserCache() (*UserCache, error) {
	path, err := GetConfigFilePath(CacheFileName)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cache UserCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}