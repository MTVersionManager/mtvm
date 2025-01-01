package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config is the application configuration
type Config struct {
	PluginDir  string `json:"pluginDir"`
	InstallDir string `json:"installDir"`
	PathDir    string `json:"pathDir"`
}

// isNotExist Checks if the error from viper.ReadInConfig is because of the configuration not existing
func isNotExist(err error) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)
	return ok
}

// GetConfig loads the configuration into a Config struct
func GetConfig() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("pluginDir", filepath.Join(configDir, "mtvm", "plugins"))
	defInstalldir, err := DefaultInstallDir()
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("installDir", defInstalldir)
	defPathDir, err := DefaultPathDir()
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("pathDir", defPathDir)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir + string(os.PathSeparator) + "mtvm")
	err = viper.ReadInConfig()
	if err != nil && !isNotExist(err) {
		return Config{}, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
