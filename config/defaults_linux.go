//go:build linux

package config

import (
	"os"
	"path"
)

func DefaultInstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".local", "share", "mtvm"), nil
}
