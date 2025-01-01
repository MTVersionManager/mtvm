package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MTVersionManager/mtvm/config"
)

// UpdateEntries updates the data of an entry if it exists, and adds an entry if it doesn't
func UpdateEntries(entry Entry) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(configDir, 0o666)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(configDir, "plugins.json"))
	var entryExists bool
	var entries []Entry
	if !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &entries)
		if err != nil {
			return err
		}
		for i, v := range entries {
			if entry.Name == v.Name {
				entryExists = true
				entries[i] = entry
				break
			}
		}
	}
	if !entryExists {
		entries = append(entries, entry)
	}
	data, err = json.MarshalIndent(entries, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configDir, "plugins.json"), data, 0o666)
}

// InstalledVersion returns the current version of a plugin that is installed.
// Returns an ErrNotFound if the version is not found.
func InstalledVersion(pluginName string) (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(configDir, "plugins.json"))
	if err != nil {
		return "", err
	}
	var entries []Entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return "", err
	}
	for _, v := range entries {
		if v.Name == pluginName {
			return v.Version, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrNotFound, pluginName)
}
