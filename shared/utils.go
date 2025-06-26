package shared

import (
	"errors"
	"github.com/spf13/afero"
	"os"
	"path/filepath"

	"github.com/MTVersionManager/mtvmplugin"
)

func IsVersionInstalled(tool, version string, fs afero.Fs) (bool, error) {
	_, err := fs.Stat(filepath.Join(Configuration.InstallDir, tool, version))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func LoadPlugin(tool string) (mtvmplugin.Plugin, error) {
	// var plugin mtvmplugin.Plugin
	// if strings.ToLower(tool) == "go" {
	//	plugin = &goplugin.Plugin{}
	// } else {
	return nil, errors.New("plugin support is not yet implemented")
	// }
	// return plugin, nil
}

func IsNotFound(err error) bool {
	var notFound NotFoundError
	return errors.As(err, &notFound)
}
