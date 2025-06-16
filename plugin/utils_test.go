package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("want error to contain ErrNotFound, got error not containing ErrNotFound")
	}
}

func TestInstalledVersionEmptyPluginFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte("[]"), fs)
	if err != nil {
		t.Fatal(err)
	}
	_, err = InstalledVersion("loremIpsum", fs)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("want error to contain ErrNotFound, got error not containing ErrNotFound")
	}
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
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("want no error when getting config dir, got %v", err)
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("want no error when reading plugins.json, got %v", err)
	}
	if string(data) != oneEntryJson {
		t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", oneEntryJson, string(data))
	}
}

func TestAddEntryWithExistingEntry(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(oneEntryJson), fs)
	if err != nil {
		t.Fatal(err)
	}
	err = UpdateEntries(Entry{
		Name:        "dolorSitAmet",
		Version:     "0.0.0",
		MetadataUrl: "https://example.com",
	}, fs)
	if err != nil {
		t.Fatalf("want no error when updating entries, got %v", err)
	}
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("want no error when getting config directory, got %v", err)
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("want no error when reading plugins.json, got %v", err)
	}
	if string(data) != twoEntryJson {
		t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", twoEntryJson, string(data))
	}
}

func TestUpdateExistingEntry(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(oneEntryJson), fs)
	if err != nil {
		t.Fatal(err)
	}
	err = UpdateEntries(Entry{
		Name:        "loremIpsum",
		Version:     "1.0.0",
		MetadataUrl: "https://example.com",
	}, fs)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("want no error when getting config directory, got %v", err)
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("want no error when reading plugins.json, got %v", err)
	}
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

func TestGetEntriesWithNoEntries(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(`[]`), fs)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := GetEntries(fs)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("want entries to be empty, got %v", entries)
	}
}

func TestGetEntriesWithEntries(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(twoEntryJson), fs)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := GetEntries(fs)
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
}

func TestRemoveExistingEntry(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(twoEntryJson), fs)
	if err != nil {
		t.Fatal(err)
	}
	err = RemoveEntry("dolorSitAmet", fs)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("want no error when getting config directory, got %v", err)
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("want no error when reading plugins.json, got %v", err)
	}
	if string(data) != oneEntryJson {
		t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", oneEntryJson, string(data))
	}
}

func TestRemoveEntryWithoutPluginsJson(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := RemoveEntry("loremIpsum", fs)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want error containing ErrNotFound, got %v", err)
	}
}

func TestRemoveEntryNonExistentEntry(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(oneEntryJson), fs)
	if err != nil {
		t.Fatal(err)
	}
	err = RemoveEntry("dolorSitAmet", fs)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want error containing ErrNotFound, got %v", err)
	}
}

func TestRemoveEntryInvalidJson(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(""), fs)
	if err != nil {
		t.Fatal(err)
	}
	err = RemoveEntry("loremIpsum", fs)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("want JSON syntax error, got %v", err)
	}
}

func TestRemoveEntryNoEntries(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte("[]"), fs)
	if err != nil {
		t.Fatal(err)
	}
	err = RemoveEntry("loremIpsum", fs)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want error containing ErrNotFound, got %v", err)
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
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want error containing ErrNotFound, got %v", err)
	}
}

func CreateAndWritePluginsJson(content []byte, fs afero.Fs) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("want no error when getting config directory, got %v", err)
	}
	err = fs.MkdirAll(configDir, 0o666)
	if err != nil {
		return fmt.Errorf("want no error when creating config directory, got %v", err)
	}
	file, err := fs.Create(filepath.Join(configDir, "plugins.json"))
	if err != nil {
		return fmt.Errorf("want no error when creating plugins.json, got %v", err)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("want no error when writing to plugins.json, got %v", err)
	}
	return nil
}
