package plugincmds

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/MTVersionManager/mtvm/components/fatalHandler"

	"github.com/MTVersionManager/mtvm/components/downloader"
	"github.com/MTVersionManager/mtvm/plugin"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/Masterminds/semver/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type installModel struct {
	downloader       downloader.Model
	pluginInfo       pluginDownloadInfo
	errorHandler     fatalHandler.Model
	versionInstalled bool
	step             int
	metadataUrl      string
	noDownload       bool
	done             bool
	fileSystem       afero.Fs
}

type pluginDownloadInfo struct {
	URL     string
	Name    string
	Version semver.Version
}

func initialInstallModel(url string) installModel {
	return installModel{
		downloader:  downloader.New(url, downloader.UseTitle("Downloading plugin metadata...")),
		metadataUrl: url,
		fileSystem:  afero.NewOsFs(),
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
			if m.step == 0 {
				cmds = append(cmds, loadMetadataCmd(m.downloader.GetDownloadedData()))
			} else {
				cmds = append(cmds, plugin.UpdateEntriesCmd(plugin.Entry{
					Name:        m.pluginInfo.Name,
					Version:     m.pluginInfo.Version.String(),
					MetadataUrl: m.metadataUrl,
				}, m.fileSystem))
			}
		case "UpdateEntries":
			m.done = true
			return m, tea.Quit
		}
	case plugin.Metadata:
		cmds = append(cmds, getPluginInfoCmd(msg))
	case pluginDownloadInfo:
		m.pluginInfo = msg
		forceFlagUsed, err := InstallCmd.Flags().GetBool("force")
		if err != nil {
			m.errorHandler, cmd = m.errorHandler.Update(err)
			cmds = append(cmds, cmd)
		}
		if m.pluginInfo.URL == "" {
			m.noDownload = true
			return m, tea.Quit
		}
		if forceFlagUsed {
			m.step++
			m.downloader = downloader.New(m.pluginInfo.URL, downloader.WriteToFs(filepath.Join(shared.Configuration.PluginDir, m.pluginInfo.Name+"."+shared.LibraryExtension), m.fileSystem), downloader.UseTitle("Downloading plugin..."))
			cmds = append(cmds, m.downloader.Init())
		} else {
			cmds = append(cmds, plugin.InstalledVersionCmd(msg.Name, m.fileSystem))
		}
	case plugin.VersionMsg:
		constraint, err := semver.NewConstraint("> " + string(msg))
		if err != nil {
			m.errorHandler, cmd = m.errorHandler.Update(err)
			cmds = append(cmds, cmd)
		}
		if !constraint.Check(&m.pluginInfo.Version) {
			m.versionInstalled = true
			return m, tea.Quit
		}
		m.step++
		m.downloader = downloader.New(m.pluginInfo.URL, downloader.WriteToFs(filepath.Join(shared.Configuration.PluginDir, m.pluginInfo.Name+".so"), m.fileSystem), downloader.UseTitle("Downloading plugin..."))
		cmds = append(cmds, m.downloader.Init())
	case plugin.NotFoundMsg:
		m.step++
		m.downloader = downloader.New(m.pluginInfo.URL, downloader.WriteToFs(filepath.Join(shared.Configuration.PluginDir, m.pluginInfo.Name+".so"), m.fileSystem), downloader.UseTitle("Downloading plugin..."))
		cmds = append(cmds, m.downloader.Init())
	}
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m installModel) View() string {
	if m.done {
		return fmt.Sprintf("%v Successfully installed version %v of the %v plugin\n", shared.CheckMark, m.pluginInfo.Version, m.pluginInfo.Name)
	}
	if m.versionInstalled {
		return fmt.Sprintf("You already have the latest version of the %v plugin installed.\nUse the --force or -f flag to reinstall it.\n", m.pluginInfo.Name)
	}
	if m.noDownload {
		return fmt.Sprintf("Sadly, the %v plugin does not provide a download for your system.\n", m.pluginInfo.Name)
	}
	return m.downloader.View() + "\n"
}

var InstallCmd = &cobra.Command{
	Use:     "install [plugin url]",
	Short:   "Install a plugin",
	Long:    `Install a plugin given a link to the plugin's metadata JSON`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"i", "in"},
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
				fatalHandler.Handle(model.errorHandler)
			} else {
				log.Fatal("Unexpected model type")
			}
		}
	},
}

func init() {
	InstallCmd.Flags().BoolP("force", "f", false, "force install a plugin")
}
