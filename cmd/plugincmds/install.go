package plugincmds

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/MTVersionManager/mtvm/components/downloader"
	"github.com/MTVersionManager/mtvm/components/fatalhandler"
	"github.com/MTVersionManager/mtvm/plugin"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/Masterminds/semver/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

type installModel struct {
	downloader       downloader.Model
	pluginInfo       pluginDownloadInfo
	errorHandler     fatalhandler.Model
	versionInstalled bool
}

type pluginDownloadInfo struct {
	URL     string
	Name    string
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
			Name:    metadata.Name,
			Version: *version,
		}
	}
}

func (m installModel) Init() tea.Cmd {
	return m.downloader.Init()
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		m.errorHandler, cmd = m.errorHandler.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, m.downloader.StopDownload()
		}
	case shared.SuccessMsg:
		switch msg {
		case "download":
			cmds = append(cmds, loadMetadataCmd(m.downloader.GetDownloadedData()))
		case "UpdateEntries":
			return m, tea.Quit
		}
	case plugin.Metadata:
		// fmt.Println(msg)
		cmds = append(cmds, getPluginInfoCmd(msg))
	case pluginDownloadInfo:
		m.pluginInfo = msg
		forceFlagUsed, err := InstallCmd.Flags().GetBool("force")
		if err != nil {
			m.errorHandler, cmd = m.errorHandler.Update(err)
			cmds = append(cmds, cmd)
		}
		if forceFlagUsed {
			cmds = append(cmds, plugin.UpdateEntriesCmd(plugin.Entry{
				Name:        m.pluginInfo.Name,
				Version:     m.pluginInfo.Version.String(),
				MetadataUrl: m.downloader.GetUrl(),
			}))
		} else {
			cmds = append(cmds, plugin.InstalledVersionCmd(msg.Name))
		}
	case plugin.VersionMsg:
		if m.pluginInfo.Version.String() == string(msg) {
			m.versionInstalled = true
			return m, tea.Quit
		}
		cmds = append(cmds, plugin.UpdateEntriesCmd(plugin.Entry{
			Name:        m.pluginInfo.Name,
			Version:     m.pluginInfo.Version.String(),
			MetadataUrl: m.downloader.GetUrl(),
		}))
	case plugin.NotFoundMsg:
		cmds = append(cmds, plugin.UpdateEntriesCmd(plugin.Entry{
			Name:        m.pluginInfo.Name,
			Version:     m.pluginInfo.Version.String(),
			MetadataUrl: m.downloader.GetUrl(),
		}))
	}
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m installModel) View() string {
	if m.pluginInfo != (pluginDownloadInfo{}) {
		if m.versionInstalled {
			return fmt.Sprintf("You already have the latest version of the %v plugin installed.\nUse the --force or -f flag to reinstall it.\n", m.pluginInfo.Name)
		}
		// fmt.Println("Finish")
		if m.pluginInfo.URL == "" {
			return "Sadly, that plugin does not provide a download for your system."
		}
		return fmt.Sprintf(`Plugin version: %v
Plugin URL: %v
`, m.pluginInfo.Version, m.pluginInfo.URL)
	}
	// fmt.Println("Download")
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

func init() {
	InstallCmd.Flags().BoolP("force", "f", false, "force install a plugin")
}
