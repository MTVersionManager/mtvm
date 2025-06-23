package plugincmds

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"

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
	tests := map[string]test{
		"existing download": {
			metadata: plugin.Metadata{
				Name:    "loremIpsum",
				Version: "0.0.0",
				Downloads: []plugin.Download{
					{
						OS:   runtime.GOOS,
						Arch: runtime.GOARCH,
						Url:  "https://example.com",
					},
				},
			},
			testFunc: func(t *testing.T, msg tea.Msg) {
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
					t.Fatalf("want error, got %T with content %v", msg, msg)
				}
			},
		},
		"invalid version": {
			metadata: plugin.Metadata{
				Name:    "loremIpsum",
				Version: "IAmAnInvalidVersion",
				Downloads: []plugin.Download{
					{
						OS:   runtime.GOOS,
						Arch: runtime.GOARCH,
						Url:  "https://example.com",
					},
				},
			},
			testFunc: func(t *testing.T, msg tea.Msg) {
				if err, ok := msg.(error); ok {
					if !errors.Is(err, semver.ErrInvalidSemVer) {
						t.Fatalf("want error to be semver.ErrInvalidSemVer, got %v", err)
					}
				} else {
					t.Fatalf("want error, got %T with content %v", msg, msg)
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

func TestInstallUpdateCancelQ(t *testing.T) {
	err := CancelTest(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPluginInstallUpdateCancelCtrlC(t *testing.T) {
	err := CancelTest(tea.KeyMsg{
		Type: tea.KeyCtrlC,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func CancelTest(keyPress tea.KeyMsg) error {
	model := initialInstallModel("https://example.com")
	_, cancel := context.WithCancel(context.Background())
	modelUpdated, _ := model.Update(downloader.DownloadStartedMsg{
		Cancel: cancel,
	})
	_, cmd := modelUpdated.Update(keyPress)
	if cmd == nil {
		return errors.New("want not nil command, got nil")
	}
	msg := cmd()
	if _, ok := msg.(downloader.DownloadCanceledMsg); !ok {
		return fmt.Errorf("expected returned command to return downloader.DownloadCanceledMsg, returned %v with type %T", msg, msg)
	}
	return nil
}

func TestPluginInstallUpdateEntriesSuccess(t *testing.T) {
	model := initialInstallModel("https://example.com")
	_, cmd := model.Update(shared.SuccessMsg("UpdateEntries"))
	if cmd == nil {
		t.Fatal("want not nil command, got nil")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("want command to return tea.QuitMsg, returned %T with content %v", msg, msg)
	}
}
