package plugincmds

import (
	"context"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/MTVersionManager/mtvm/components/downloader"
	"github.com/MTVersionManager/mtvm/plugin"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/Masterminds/semver/v3"
	tea "github.com/charmbracelet/bubbletea"
)

func TestGetPluginInfo(t *testing.T) {
	type test struct {
		metadata plugin.Metadata
		testFunc func(t *testing.T, msg tea.Msg)
	}
	exampleUsableDownloads := []plugin.Download{
		{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
			Url:  "https://example.com",
		},
	}
	tests := map[string]test{
		"existing download": {
			metadata: plugin.Metadata{
				Name:      "loremIpsum",
				Version:   "0.0.0",
				Downloads: exampleUsableDownloads,
			},
			testFunc: func(t *testing.T, msg tea.Msg) {
				if downloadInfo, ok := msg.(pluginDownloadInfo); ok {
					assert.Equalf(t, "loremIpsum", downloadInfo.Name, "want name to be 'loremIpsum', got name '%v'", downloadInfo.Name)
					assert.Equalf(t, "https://example.com", downloadInfo.Url, "want url to be 'https://example.com', got url '%v'", downloadInfo.Url)
					compareVersionTo := semver.New(0, 0, 0, "", "")
					if !compareVersionTo.Equal(downloadInfo.Version) {
						t.Errorf("Want version 0.0.0 got %v", downloadInfo.Version.String())
					}
				} else if err, ok := msg.(error); ok {
					t.Errorf("want no error, got %v", err)
				} else {
					t.Errorf("want pluginDownloadInfo returned, got %T with content %v", msg, msg)
				}
			},
		},
		"no download": {
			metadata: plugin.Metadata{
				Name:    "loremIpsum",
				Version: "0.0.0",
				Downloads: []plugin.Download{
					{
						OS: func() string {
							if runtime.GOOS == "imaginaryOS" {
								return "fakeOS"
							}
							return "imaginaryOS"
						}(),
						Arch: func() string {
							if runtime.GOARCH == "imaginaryArch" {
								return "fakeArch"
							}
							return "imaginaryArch"
						}(),
						Url: "https://example.com",
					},
				},
			},
			testFunc: func(t *testing.T, msg tea.Msg) {
				if err, ok := msg.(error); ok {
					shared.AssertIsNotFoundError(t, err, "download", shared.Source{
						File:     "cmd/plugincmds/install.go",
						Function: "getPluginInfoCmd(metadata plugin.Metadata) tea.Cmd",
					})
				} else {
					t.Errorf("want error, got %T with content %v", msg, msg)
				}
			},
		},
		"invalid version": {
			metadata: plugin.Metadata{
				Name:      "loremIpsum",
				Version:   "IAmAnInvalidVersion",
				Downloads: exampleUsableDownloads,
			},
			testFunc: func(t *testing.T, msg tea.Msg) {
				if err, ok := msg.(error); ok {
					assert.ErrorIs(t, err, semver.ErrInvalidSemVer)
				} else {
					t.Errorf("want error, got %T with content %v", msg, msg)
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			msg := getPluginInfoCmd(tt.metadata)()
			tt.testFunc(t, msg)
		})
	}
}

func TestInstallUpdateCancel(t *testing.T) {
	tests := map[string]tea.KeyMsg{
		"ctrl+c": {
			Type: tea.KeyCtrlC,
		},
		"q": {
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		},
	}
	for name, keyPress := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			model := initialInstallModel("https://example.com")
			_, cancel := context.WithCancel(context.Background())
			modelUpdated, _ := model.Update(downloader.DownloadStartedMsg{
				Cancel: cancel,
			})
			_, cmd := modelUpdated.Update(keyPress)
			require.NotNil(t, cmd, "want not nil command, got nil")
			msg := cmd()
			assert.IsType(t, downloader.DownloadCanceledMsg{}, msg)
		})
	}
}
