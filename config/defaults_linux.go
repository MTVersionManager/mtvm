//go:build linux

package config

import (
	"os"
	"path/filepath"
)

func DefaultInstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "mtvm"), nil
}

func DefaultPathDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "bin", "mtvm"), nil
}

func GetConfigDir() (string, error) {
	envDir := os.Getenv("MTVM_CONFIG_DIR")
	if envDir == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(userConfigDir, "mtvm"), nil
	}
	return envDir, nil
}
