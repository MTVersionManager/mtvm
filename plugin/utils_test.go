package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MTVersionManager/mtvm/config"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/spf13/afero"
)

var oneEntryJson = `[
	{
		"name": "loremIpsum",
		"version": "0.0.0",
		"metadataUrl": "https://example.com"
	}
]`

var twoEntryJson = `[
	{
		"name": "loremIpsum",
		"version": "0.0.0",
		"metadataUrl": "https://example.com"
	},
	{
		"name": "dolorSitAmet",
		"version": "0.0.0",
		"metadataUrl": "https://example.com"
	}
]`

func TestInstalledVersionNoPluginsJson(t *testing.T) {
	_, err := InstalledVersion("loremIpsum", afero.NewMemMapFs())
	shared.AssertIsNotFoundError(t, err, "plugins.json", shared.Source{
		File:     "plugin/utils.go",
		Function: "InstalledVersion(pluginName string, fs afero.Fs) (string, error)",
	})
}

func TestInstalledVersionWithPluginsJson(t *testing.T) {
	testFuncErrNotFound := func(t *testing.T, _ string, err error) {
		shared.AssertIsNotFoundError(t, err, "entry", shared.Source{
			File:     "plugin/utils.go",
			Function: "InstalledVersion(pluginName string, fs afero.Fs) (string, error)",
		})
	}
	tests := map[string]struct {
		pluginsJsonContent []byte
		pluginName         string
		testFunc           func(t *testing.T, version string, err error)
	}{
		"empty plugins.json": {
			pluginsJsonContent: []byte(`[]`),
			pluginName:         "loremIpsum",
			testFunc:           testFuncErrNotFound,
		},
		"non-existent entry": {
			pluginsJsonContent: []byte(oneEntryJson),
			pluginName:         "dolorSitAmet",
			testFunc:           testFuncErrNotFound,
		},
		"invalid json": {
			pluginsJsonContent: []byte(""),
			pluginName:         "loremIpsum",
			testFunc: func(t *testing.T, version string, err error) {
				checkIfJsonSyntaxError(t, err)
				assert.Emptyf(t, version, "want version to be empty, got %v", version)
			},
		},
		"existing entry": {
			pluginsJsonContent: []byte(oneEntryJson),
			pluginName:         "loremIpsum",
			testFunc: func(t *testing.T, version string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "0.0.0", version)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, tt.pluginsJsonContent, fs)
			version, err := InstalledVersion(tt.pluginName, fs)
			tt.testFunc(t, version, err)
		})
	}
}

func TestAddFirstEntryNoPluginsJson(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := UpdateEntries(Entry{
		Name:        "loremIpsum",
		Version:     "0.0.0",
		MetadataUrl: "https://example.com",
	}, fs)
	assert.NoError(t, err)
	data := readPluginsJson(t, fs)
	assert.Equal(t, oneEntryJson, string(data))
}

func TestUpdateEntryWithPluginsJson(t *testing.T) {
	tests := map[string]struct {
		pluginsJsonContent []byte
		entry              Entry
		wantsError         bool
		testFunc           func(t *testing.T, fs afero.Fs, err error)
	}{
		"add second": {
			pluginsJsonContent: []byte(oneEntryJson),
			entry: Entry{
				Name:        "dolorSitAmet",
				Version:     "0.0.0",
				MetadataUrl: "https://example.com",
			},
			wantsError: false,
			testFunc: func(t *testing.T, fs afero.Fs, err error) {
				data := readPluginsJson(t, fs)
				assert.Equal(t, twoEntryJson, string(data))
			},
		},
		"update existing": {
			pluginsJsonContent: []byte(oneEntryJson),
			entry: Entry{
				Name:        "loremIpsum",
				Version:     "1.0.0",
				MetadataUrl: "https://example.com",
			},
			wantsError: false,
			testFunc: func(t *testing.T, fs afero.Fs, err error) {
				data := readPluginsJson(t, fs)
				expected := `[
	{
		"name": "loremIpsum",
		"version": "1.0.0",
		"metadataUrl": "https://example.com"
	}
]`
				assert.Equal(t, expected, string(data))
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, tt.pluginsJsonContent, fs)
			err := UpdateEntries(tt.entry, fs)
			if tt.wantsError {
				assert.Error(t, err)
			}
			if !tt.wantsError {
				assert.NoError(t, err)
			}
			tt.testFunc(t, fs, err)
		})
	}
}

func TestGetEntriesWithNoPluginsJson(t *testing.T) {
	fs := afero.NewMemMapFs()
	entries, err := GetEntries(fs)
	assert.NoError(t, err)
	assert.Nilf(t, entries, "want entries to be nil, got %v", err)
}

func TestGetEntriesWithPluginsJson(t *testing.T) {
	tests := map[string]struct {
		pluginsJsonContent []byte
		testFunc           func(t *testing.T, entries []Entry, err error)
	}{
		"empty": {
			pluginsJsonContent: []byte(`[]`),
			testFunc: func(t *testing.T, entries []Entry, err error) {
				assert.NoError(t, err)
				assert.Lenf(t, entries, 0, "want entries to be empty, got %v", entries)
			},
		},
		"two entries": {
			pluginsJsonContent: []byte(twoEntryJson),
			testFunc: func(t *testing.T, entries []Entry, err error) {
				assert.NoError(t, err)
				assert.Len(t, entries, 2, "want 2 entries")
				assert.Equalf(t, "loremIpsum", entries[0].Name, "want first entry name to be 'loremIpsum', got %v", entries[0].Name)
				assert.Equalf(t, "dolorSitAmet", entries[1].Name, "want second entry name to be 'dolorSitAmet', got %v", entries[1].Name)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, tt.pluginsJsonContent, fs)
			entries, err := GetEntries(fs)
			tt.testFunc(t, entries, err)
		})
	}
}

func TestRemoveEntryWithoutPluginsJson(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := RemoveEntry("loremIpsum", fs)
	shared.AssertIsNotFoundError(t, err, "plugins.json", shared.Source{
		File:     "plugin/utils.go",
		Function: "RemoveEntry(pluginName string, fs afero.Fs) error",
	})
}

func TestRemoveEntryWithPluginsJson(t *testing.T) {
	testFuncErrNotFound := func(t *testing.T, _ afero.Fs, err error) {
		shared.AssertIsNotFoundError(t, err, "entry", shared.Source{
			File:     "plugin/utils.go",
			Function: "RemoveEntry(pluginName string, fs afero.Fs) error",
		})
	}
	tests := map[string]struct {
		pluginToRemove     string
		pluginsJsonContent []byte
		testFunc           func(t *testing.T, fs afero.Fs, err error)
	}{
		"existing entry": {
			pluginToRemove:     "dolorSitAmet",
			pluginsJsonContent: []byte(twoEntryJson),
			testFunc: func(t *testing.T, fs afero.Fs, err error) {
				assert.NoError(t, err)
				data := readPluginsJson(t, fs)
				assert.Equal(t, oneEntryJson, string(data))
			},
		},
		"non-existent entry": {
			pluginToRemove:     "dolorSitAmet",
			pluginsJsonContent: []byte(oneEntryJson),
			testFunc:           testFuncErrNotFound,
		},
		"no entries": {
			pluginToRemove:     "loremIpsum",
			pluginsJsonContent: []byte(`[]`),
			testFunc:           testFuncErrNotFound,
		},
		"invalid json": {
			pluginToRemove:     "loremIpsum",
			pluginsJsonContent: []byte(""),
			testFunc: func(t *testing.T, _ afero.Fs, err error) {
				checkIfJsonSyntaxError(t, err)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, tt.pluginsJsonContent, fs)
			err := RemoveEntry(tt.pluginToRemove, fs)
			tt.testFunc(t, fs, err)
		})
	}
}

func TestRemoveExisting(t *testing.T) {
	fs := afero.NewMemMapFs()
	var err error
	shared.Configuration, err = config.GetConfig()
	require.NoError(t, err, "when getting configuration")
	err = fs.MkdirAll(shared.Configuration.PluginDir, 0o777)
	require.NoError(t, err, "when creating plugin directory")
	pluginPath := filepath.Join(shared.Configuration.PluginDir, "loremIpsum"+shared.LibraryExtension)
	_, err = fs.Create(pluginPath)
	require.NoError(t, err, "when creating plugin file")
	err = Remove("loremIpsum", fs)
	require.NoError(t, err)
	_, err = fs.Stat(pluginPath)
	assert.Error(t, err, "when statting plugin file")
	if !os.IsNotExist(err) {
		t.Errorf("want file does not exist error, got %v (stat)", err)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := Remove("loremIpsum", fs)
	shared.AssertIsNotFoundError(t, err, "plugin", shared.Source{
		File:     "plugin/utils.go",
		Function: "Remove(pluginName string, fs afero.Fs) error",
	})
}

func createAndWritePluginsJson(t *testing.T, content []byte, fs afero.Fs) {
	configDir, err := config.GetConfigDir()
	require.NoError(t, err, "when getting config directory")
	err = fs.MkdirAll(configDir, 0o666)
	require.NoError(t, err, "when creating config directory")
	file, err := fs.Create(filepath.Join(configDir, "plugins.json"))
	require.NoError(t, err, "when creating plugins.json")
	defer func() {
		err := file.Close()
		assert.NoError(t, err, "when closing plugins.json")
	}()
	_, err = file.Write(content)
	require.NoError(t, err, "when writing to plugins.json")
}

func readPluginsJson(t *testing.T, fs afero.Fs) []byte {
	configDir, err := config.GetConfigDir()
	require.NoError(t, err, "when getting config directory")
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	require.NoError(t, err, "when reading plugins.json")
	return data
}

func checkIfJsonSyntaxError(t *testing.T, err error) {
	require.Error(t, err)
	var syntaxError *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxError)
}
