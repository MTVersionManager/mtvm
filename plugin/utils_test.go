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

func TestInstalledVersionNoPluginFile(t *testing.T) {
	_, err := InstalledVersion("loremIpsum", afero.NewMemMapFs())
	checkIfErrNotFound(t, err)
}

func TestInstalledVersionEmptyPluginFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	createAndWritePluginsJson(t, []byte("[]"), fs)
	_, err := InstalledVersion("loremIpsum", fs)
	checkIfErrNotFound(t, err)
}

func TestAddFirstEntryNoPluginFile(t *testing.T) {
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
	tests := []struct {
		name               string
		pluginsJsonContent []byte
		entry              Entry
		wantsError         bool
		testFunc           func(t *testing.T, fs afero.Fs, err error)
	}{
		{
			name:               "AddSecondEntry",
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
		{
			name:               "UpdateExistingEntry",
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, test.pluginsJsonContent, fs)
			err := UpdateEntries(test.entry, fs)
			if test.wantsError && err == nil {
				t.Fatal("want error, got nil")
			}
			if !test.wantsError && err != nil {
				t.Fatalf("want no error, got %v", err)
			}
			test.testFunc(t, fs, err)
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
	tests := []struct {
		name               string
		pluginsJsonContent []byte
		testFunc           func(t *testing.T, entries []Entry, err error)
	}{
		{
			name:               "NoEntries",
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
		{
			name:               "TwoEntries",
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, test.pluginsJsonContent, fs)
			entries, err := GetEntries(fs)
			test.testFunc(t, entries, err)
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
	tests := []struct {
		name               string
		pluginToRemove     string
		pluginsJsonContent []byte
		testFunc           func(t *testing.T, fs afero.Fs, err error)
	}{
		{
			name:               "ExistingEntry",
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
		{
			name:               "NonExistentEntry",
			pluginToRemove:     "dolorSitAmet",
			pluginsJsonContent: []byte(oneEntryJson),
			testFunc:           testFuncErrNotFound,
		},
		{
			name:               "NoEntries",
			pluginToRemove:     "loremIpsum",
			pluginsJsonContent: []byte(`[]`),
			testFunc:           testFuncErrNotFound,
		},
		{
			name:               "InvalidJson",
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			createAndWritePluginsJson(t, test.pluginsJsonContent, fs)
			err := RemoveEntry(test.pluginToRemove, fs)
			test.testFunc(t, fs, err)
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
