package shared

import (
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