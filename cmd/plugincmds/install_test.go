package plugincmds

import (
	"runtime"
	"testing"

	"github.com/MTVersionManager/mtvm/plugin"
	"github.com/Masterminds/semver/v3"
)

func TestGetPluginInfo(t *testing.T) {
	msg := getPluginInfoCmd(plugin.Metadata{
		Name:    "loremIpsum",
		Version: "0.0.0",
		Downloads: []plugin.Download{
			{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
				Url:  "https://example.com",
			},
		},
	})()
	if downloadInfo, ok := msg.(pluginDownloadInfo); ok {
		if downloadInfo.Name != "loremIpsum" {
			t.Fatalf("want name to be 'loremIpsum', got name '%v'", downloadInfo.Name)
		}
		if downloadInfo.Url != "https://example.com" {
			t.Fatalf("want url to be 'https://example.com', got url '%v'", downloadInfo.Url)
		}
		compareVersionTo := semver.New(0, 0, 0, "", "")
		if !compareVersionTo.Equal(downloadInfo.Version) {
			t.Fatalf("Want version 0.0.0 got %v", downloadInfo.Version.String())
		}
	} else if err, ok := msg.(error); ok {
		t.Fatalf("want no error, got %v", err)
	} else {
		t.Fatalf("want pluginDownloadInfo returned, got %T with content %v", msg, msg)
	}
}
