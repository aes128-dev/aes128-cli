package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigDirName = "aes128-cli"
	TokenFileName = "session.token"
)

func getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appConfigDir := filepath.Join(configDir, ConfigDirName)
	if err := os.MkdirAll(appConfigDir, 0750); err != nil {
		return "", err
	}
	return appConfigDir, nil
}

func SaveToken(token string) error {
	dir, err := getConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, TokenFileName)
	return os.WriteFile(path, []byte(token), 0600)
}

func ReadToken() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, TokenFileName)
	token, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(token), nil
}