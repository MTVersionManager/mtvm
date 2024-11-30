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

var CheckMark string = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).SetString("✓").String()

func IsVersionInstalled(tool string, version string) (bool, error) {
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
