package config

import (
	"github.com/spf13/viper"
	"os"
)

// Config is the application configuration
type Config struct {
	PluginDir  string
	InstallDir string
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
