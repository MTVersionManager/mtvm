package shared

import (
	"errors"
	"github.com/MTVersionManager/goplugin"
	"github.com/MTVersionManager/mtvmplugin"
	"os"
	"path"

	"github.com/MTVersionManager/mtvm/config"
)

var Configuration config.Config

func IsVersionInstalled(tool string, version string) (bool, error) {
	_, err := os.Stat(path.Join(Configuration.InstallDir, tool, version))
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
	if tool == "go" {
		plugin = &goplugin.Plugin{}
	} else {
		return nil, errors.New("plugin support is not yet implemented")
	}
	return plugin, nil
}
