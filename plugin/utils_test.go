package plugin

import (
	"errors"
	"testing"

	"github.com/spf13/afero"
	//"os"
	//"path/filepath"
	//"github.com/MTVersionManager/mtvm/config"
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

// TODO: Finish this test after doing afero
//func TestAddFirstEntry(t *testing.T) {
//	err := UpdateEntries(Entry{
//		Name: "loremIpsum",
//		Version: "0.0.0",
//		MetadataUrl: "http://example.com",
//	})
//	if err != nil {
//		t.Fatalf("Want no error, got %v", err)
//	}
//	configDir, err := config.GetConfigDir()
//	if err != nil {
//		t.Fatalf("Want no error when getting config dir, got %v", err)
//	}
//	data, err := os.ReadFile(filepath.Join(configDir, "plugins.json"))
//	if err != nil {
//		t.Fatalf("Want no error when reading plugins.json, got %v", err)
//	}
//	t.Log(string(data))
//}
