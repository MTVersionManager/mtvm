package plugin

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/MTVersionManager/mtvm/config"
	"github.com/spf13/afero"
)

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
	t.Log(err)
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("want error to contain ErrNotFound, got error not containing ErrNotFound")
	}
}

func TestAddFirstEntryNoPluginFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := UpdateEntries(Entry{
		Name:        "loremIpsum",
		Version:     "0.0.0",
		MetadataUrl: "http://example.com",
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
	expected := `[
	{
		"name": "loremIpsum",
		"version": "0.0.0",
		"metadataUrl": "https://example.com"
	}
]`
	if string(data) != expected {
		t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", expected, string(data))
	}
}

func TestAddEntryWithExistingEntry(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := CreateAndWritePluginsJson([]byte(`[
	{
		"name": "loremIpsum",
		"version": "0.0.0",
		"metadataUrl": "https://example.com"
	}
]`), fs)
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
	expected := `[
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
	if string(data) != expected {
		t.Fatalf("want plugins.json to contain\n%v\ngot plugins.json containing\n%v", expected, string(data))
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
