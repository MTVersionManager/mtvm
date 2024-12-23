package shared

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/MTVersionManager/goplugin"
	"github.com/MTVersionManager/mtvmplugin"
	"github.com/charmbracelet/lipgloss"

	"github.com/MTVersionManager/mtvm/config"
)

var Configuration config.Config

type SuccessMsg string

var CheckMark = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).SetString("✓").String()

type PluginMetadata struct {
	Name      string `json:"name" validate:"required"`
	Version   string `json:"version" validate:"required,semver"`
	Downloads []struct {
		OS       string `json:"os" validate:"required"`
		Arch     string `json:"arch" validate:"required"`
		URL      string `json:"url" validate:"required,http_url"`
		Checksum string `json:"checksum"`
	} `json:"downloads" validate:"required,dive"`
}

func IsVersionInstalled(tool, version string) (bool, error) {
	_, err := os.Stat(filepath.Join(Configuration.InstallDir, tool, version))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func LoadPlugin(tool string) (mtvmplugin.Plugin, error) {
	var plugin mtvmplugin.Plugin
	if strings.ToLower(tool) == "go" {
		plugin = &goplugin.Plugin{}
	} else {
		return nil, errors.New("plugin support is not yet implemented")
	}
	return plugin, nil
}
