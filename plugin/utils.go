package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MTVersionManager/mtvm/shared"
	"github.com/spf13/afero"

	"github.com/MTVersionManager/mtvm/config"
)

// UpdateEntries updates the data of an entry if it exists and adds an entry if it doesn't
func UpdateEntries(entry Entry, fs afero.Fs) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	err = fs.MkdirAll(configDir, 0o666)
	if err != nil {
		return err
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
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
	return afero.WriteFile(fs, filepath.Join(configDir, "plugins.json"), data, 0o666)
}

// InstalledVersion returns the current version of a plugin that is installed.
// Returns a NotFoundError if the version is not found.
func InstalledVersion(pluginName string, fs afero.Fs) (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", shared.NotFoundError{
				Thing: "plugins.json",
				Source: shared.Source{
					File:     "plugin/utils.go",
					Function: "InstalledVersion(pluginName string, fs afero.Fs) (string, error)",
				},
			}
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
	return "", fmt.Errorf("%w: %s", shared.NotFoundError{
		Thing: "entry",
		Source: shared.Source{
			File:     "plugin/utils.go",
			Function: "InstalledVersion(pluginName string, fs afero.Fs) (string, error)",
		},
	}, pluginName)
}

// GetEntries returns a list of installed plugins and an error
// and returns a nil slice if plugins.json is not present or if there is an error
func GetEntries(fs afero.Fs) ([]Entry, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
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

func RemoveEntry(pluginName string, fs afero.Fs) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return shared.NotFoundError{
				Thing: "plugins.json",
				Source: shared.Source{
					File:     "plugin/utils.go",
					Function: "RemoveEntry(pluginName string, fs afero.Fs) error",
				},
			}
		}
		return err
	}
	var entries []Entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return err
	}
	notFound := shared.NotFoundError{
		Thing: "entry",
		Source: shared.Source{
			File:     "plugin/utils.go",
			Function: "RemoveEntry(pluginName string, fs afero.Fs) error",
		},
	}
	if len(entries) == 0 {
		return notFound
	}
	removed := make([]Entry, 0, len(entries)-1)
	for _, v := range entries {
		if v.Name != pluginName {
			removed = append(removed, v)
		}
	}
	if len(removed) == len(entries) {
		return notFound
	}
	data, err = json.MarshalIndent(removed, "", "	")
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, filepath.Join(configDir, "plugins.json"), data, 0o666)
}

func Remove(pluginName string, fs afero.Fs) error {
	err := fs.Remove(filepath.Join(shared.Configuration.PluginDir, pluginName+shared.LibraryExtension))
	if os.IsNotExist(err) {
		return shared.NotFoundError{
			Thing: "plugin",
			Source: shared.Source{
				File:     "plugin/utils.go",
				Function: "Remove(pluginName string, fs afero.Fs) error",
			},
		}
	}
	return err
}
