package plugincmds

import (
	"encoding/json"
	"fmt"
	"github.com/MTVersionManager/mtvm/components/downloader"
	"github.com/MTVersionManager/mtvm/components/fatalhandler"
	"github.com/MTVersionManager/mtvm/plugin"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/Masterminds/semver/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime"
)

type installModel struct {
	downloader   downloader.Model
	pluginInfo   pluginDownloadInfo
	errorHandler fatalhandler.Model
}

type pluginDownloadInfo struct {
	URL     string
	Version semver.Version
}

func initialInstallModel(url string) installModel {
	return installModel{
		downloader: downloader.New(url, downloader.UseTitle("Downloading plugin metadata...")),
	}
}

func loadMetadata(rawData []byte) (plugin.Metadata, error) {
	var metadata plugin.Metadata
	err := json.Unmarshal(rawData, &metadata)
	if err != nil {
		return metadata, err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}

func loadMetadataCmd(rawData []byte) tea.Cmd {
	return func() tea.Msg {
		metadata, err := loadMetadata(rawData)
		if err != nil {
			return err
		}
		return metadata
	}
}

func getPluginInfoCmd(metadata plugin.Metadata) tea.Cmd {
	return func() tea.Msg {
		version, err := semver.NewVersion(metadata.Version)
		if err != nil {
			return err
		}
		var url string
		for _, v := range metadata.Downloads {
			if v.OS == runtime.GOOS && v.Arch == runtime.GOARCH {
				url = v.URL
			}
		}
		return pluginDownloadInfo{
			URL:     url,
			Version: *version,
		}
	}
}

func (m installModel) Init() tea.Cmd {
	return m.downloader.Init()
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, m.downloader.StopDownload()
		}
	case shared.SuccessMsg:
		if msg == "download" {
			cmds = append(cmds, loadMetadataCmd(m.downloader.GetDownloadedData()))
		}
	case plugin.Metadata:
		//fmt.Println(msg)
		cmds = append(cmds, getPluginInfoCmd(msg))
	case pluginDownloadInfo:
		m.pluginInfo = msg
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.errorHandler, cmd = m.errorHandler.Update(msg)
	cmds = append(cmds, cmd)
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m installModel) View() string {
	if m.pluginInfo != (pluginDownloadInfo{}) {
		//fmt.Println("Finish")
		if m.pluginInfo.URL == "" {
			return "Sadly, that plugin does not provide a download for your system."
		}
		return fmt.Sprintf("Plugin version: %v\nPlugin URL: %v\n", m.pluginInfo.Version, m.pluginInfo.URL)
	}
	//fmt.Println("Download")
	return m.downloader.View() + "\n"
}

var InstallCmd = &cobra.Command{
	Use:   "install [plugin url]",
	Short: "Install a plugin",
	Long:  `Install a plugin given a link to the plugin's metadata JSON`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Var(args[0], "http_url")
		if err != nil {
			fmt.Println("Please enter a valid http url")
			os.Exit(1)
		}
		p := tea.NewProgram(initialInstallModel(args[0]))
		if model, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
		} else {
			if model, ok := model.(installModel); ok {
				fatalhandler.Handle(model.errorHandler)
			} else {
				log.Fatal("Unexpected model type")
			}
		}
	},
}
