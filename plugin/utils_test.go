package plugin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/MTVersionManager/mtvm/config"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/spf13/afero"
)

var oneEntryJson string = `[
	{
		"name": "loremIpsum",
		"version": "0.0.0",
		"metadataUrl": "https://example.com"
	}
]`

var twoEntryJson string = `[
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
	checkIfErrNotFound(t, err)
}

func TestInstalledVersionWithPluginsJson(t *testing.T) {
	testFuncErrNotFound := func(t *testing.T, _ string, err error) {
		checkIfErrNotFound(t, err)
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
				if err == nil {
					t.Fatal("want error, got nil")
				}
				if _, ok := err.(*json.SyntaxError); !ok {
					t.Fatalf("want JSON syntax error, got %v", err)
				}
				if version != "" {
					t.Fatalf("want version to be empty, got %v", version)
				}
			},
		},
		"existing entry": {
			pluginsJsonContent: []byte(oneEntryJson),
			pluginName:         "loremIpsum",
			testFunc: func(t *testing.T, version string, err error) {
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
				if version != "0.0.0" {
					t.Fatalf("want version to be '0.0.0', got %v", version)
				}
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
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	data := readPluginsJson(t, fs)
	if string(data) != oneEntryJson {
		t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", oneEntryJson, string(data))
	}
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
				if string(data) != twoEntryJson {
					t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", twoEntryJson, string(data))
				}
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
				if string(data) != expected {
					t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", expected, string(data))
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, tt.pluginsJsonContent, fs)
			err := UpdateEntries(tt.entry, fs)
			if tt.wantsError && err == nil {
				t.Fatal("want error, got nil")
			}
			if !tt.wantsError && err != nil {
				t.Fatalf("want no error, got %v", err)
			}
			tt.testFunc(t, fs, err)
		})
	}
}

func TestGetEntriesWithNoPluginsJson(t *testing.T) {
	fs := afero.NewMemMapFs()
	entries, err := GetEntries(fs)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	if entries != nil {
		t.Fatalf("want entries to be nil, got %v", entries)
	}
}

func TestGetEntriesWithPluginsJson(t *testing.T) {
	tests := map[string]struct {
		pluginsJsonContent []byte
		testFunc           func(t *testing.T, entries []Entry, err error)
	}{
		"empty": {
			pluginsJsonContent: []byte(`[]`),
			testFunc: func(t *testing.T, entries []Entry, err error) {
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
				if len(entries) != 0 {
					t.Fatalf("want entries to be empty, got %v", entries)
				}
			},
		},
		"two entries": {
			pluginsJsonContent: []byte(twoEntryJson),
			testFunc: func(t *testing.T, entries []Entry, err error) {
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
				if len(entries) != 2 {
					t.Fatalf("want 2 entries, got %v entries containing %v", len(entries), entries)
				}
				if entries[0].Name != "loremIpsum" {
					t.Fatalf("wanted first entry name to be 'loremIpsum', got %v", entries[0].Name)
				}
				if entries[1].Name != "dolorSitAmet" {
					t.Fatalf("wanted second entry name to be 'dolorSitAmet', got %v", entries[1].Name)
				}
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
	checkIfErrNotFound(t, err)
}

func TestRemoveEntryWithPluginsJson(t *testing.T) {
	testFuncErrNotFound := func(t *testing.T, _ afero.Fs, err error) {
		checkIfErrNotFound(t, err)
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
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
				data := readPluginsJson(t, fs)
				if string(data) != oneEntryJson {
					t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", oneEntryJson, string(data))
				}
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
				if err == nil {
					t.Fatal("want error, got nil")
				}
				if _, ok := err.(*json.SyntaxError); !ok {
					t.Fatalf("want JSON syntax error, got %v", err)
				}
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
	if err != nil {
		t.Fatalf("want no error when getting configuration, got %v", err)
	}
	err = fs.MkdirAll(shared.Configuration.PluginDir, 0o777)
	if err != nil {
		t.Fatalf("want no error when creating plugin directory, got %v", err)
	}
	pluginPath := filepath.Join(shared.Configuration.PluginDir, "loremIpsum"+shared.LibraryExtension)
	_, err = fs.Create(pluginPath)
	if err != nil {
		t.Fatalf("want no error when creating plugin file, got %v", err)
	}
	err = Remove("loremIpsum", fs)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	_, err = fs.Stat(pluginPath)
	if err == nil {
		t.Fatal("want error, got nil (stat)")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("want file does not exist error, got %v (stat)", err)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := Remove("loremIpsum", fs)
	checkIfErrNotFound(t, err)
}

func createAndWritePluginsJson(t *testing.T, content []byte, fs afero.Fs) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("want no error when getting config directory, got %v", err)
	}
	err = fs.MkdirAll(configDir, 0o666)
	if err != nil {
		t.Fatalf("want no error when creating config directory, got %v", err)
	}
	file, err := fs.Create(filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("want no error when creating plugins.json, got %v", err)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("want no error when writing to plugins.json, got %v", err)
	}
}

func readPluginsJson(t *testing.T, fs afero.Fs) []byte {
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("want no error when getting config directory, got %v", err)
		return nil
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("want no error when reading plugins.json, got %v", err)
	}
	return data
}

func checkIfErrNotFound(t *testing.T, err error) {
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want error to contain ErrNotFound, got error not containing ErrNotFound")
	}
}
