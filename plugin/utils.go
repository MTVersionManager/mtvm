package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/MTVersionManager/mtvm/shared"
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
			if entry.Name == v.Name && entry.MetadataUrl == v.MetadataUrl {
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
		if os.IsNotExist(err) {
			return "", ErrNotFound
		}
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

// GetEntries returns a list of installed plugins and an error
// and returns a nil slice if plugins.json is not present or if there is an error
func GetEntries() ([]Entry, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(configDir, "plugins.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
	}
	var entries []Entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func RemoveEntry(pluginName string) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(configDir, "plugins.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	var entries []Entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return err
	}
	removed := make([]Entry, 0, len(entries)-1)
	for _, v := range entries {
		if v.Name != pluginName {
			removed = append(removed, v)
		}
	}
	if len(removed) == len(entries) {
		return ErrNotFound
	}
	data, err = json.MarshalIndent(removed, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configDir, "plugins.json"), data, 0o666)
}

func RemovePlugin(pluginName string) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(configDir, "plugins", pluginName+shared.LibraryExtension))
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}
