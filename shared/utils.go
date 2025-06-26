package shared

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

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

func LoadPlugin(_ string) (mtvmplugin.Plugin, error) {
	return nil, errors.New("plugin support is not yet implemented")
}

func IsNotFound(err error) bool {
	return errors.As(err, &NotFoundError{})
}
