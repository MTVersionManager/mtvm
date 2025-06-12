package plugin

import (
	"errors"
	"testing"

	"path/filepath"

	"github.com/MTVersionManager/mtvm/config"
	"github.com/spf13/afero"
)

func TestInstalledVersionNoPluginFile(t *testing.T) {
	_, err := InstalledVersion("loremIpsum", afero.NewMemMapFs())
	if err == nil {
		t.Fatal("Want error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("Want error to contain ErrNotFound, got error not containing ErrNotFound")
	}
}

func TestInstalledVersionEmptyPluginFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("Want no error when getting config directory, got %v", err)
	}
	err = fs.MkdirAll(configDir, 0o666)
	if err != nil {
		t.Fatalf("Want no error when creating config directory, got %v", err)
	}
	file, err := fs.Create(filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("Want no error when creating plugins.json, got %v", err)
	}
	defer file.Close()
	_, err = file.Write([]byte("[]"))
	if err != nil {
		t.Fatalf("Want no error when writing to plugins.json, got %v", err)
	}
	_, err = InstalledVersion("loremIpsum", fs)
	if err == nil {
		t.Fatal("Want error, got nil")
	}
	t.Log(err)
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("Want error to contain ErrNotFound, got error not containing ErrNotFound")
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
		t.Fatalf("Want no error, got %v", err)
	}
	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("Want no error when getting config dir, got %v", err)
	}
	data, err := afero.ReadFile(fs, filepath.Join(configDir, "plugins.json"))
	if err != nil {
		t.Fatalf("Want no error when reading plugins.json, got %v", err)
	}
	expected := `[
	{
		"name": "loremIpsum",
		"version": "0.0.0",
		"metadataUrl": "http://example.com"
	}
]`
	if string(data) != expected {
		t.Fatalf("Want\n%v\ngot\n%v", expected, string(data))
	}
}
