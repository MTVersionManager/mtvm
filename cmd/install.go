package cmd

import (
	"fmt"
	"github.com/MTVersionManager/goplugin"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/MTVersionManager/mtvmplugin"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

type installModel struct {
	progressChannel chan float64
	plugin          mtvmplugin.Plugin
	installing      bool
	progress        float64
	progressBar     progress.Model
	version         string
	spinner         spinner.Model
	pluginName      string
}

type progressMsg float64
type downloadedMsg bool
type installedMsg bool

func installInitialModel(plugin mtvmplugin.Plugin, pluginName string, version string) installModel {
	progressChannel := make(chan float64)
	progressBar := progress.New(progress.WithDefaultGradient())
	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	return installModel{
		progressChannel: progressChannel,
		progressBar:     progressBar,
		version:         version,
		plugin:          plugin,
		installing:      true,
		spinner:         spinnerModel,
		pluginName:      pluginName,
	}
}

func (m installModel) Init() tea.Cmd {
	return tea.Batch(download(m.plugin, m.version, m.progressChannel), waitForProgress(m.progressChannel), m.spinner.Tick)
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case error:
		log.Fatal(msg)
	case downloadedMsg:
		m.progress = 1
		m.installing = false
		return m, install(m.plugin, shared.Configuration.InstallDir, m.pluginName, m.version)
	case installedMsg:
		return m, tea.Quit
	case progressMsg:
		m.progress = float64(msg)
		cmds = append(cmds, waitForProgress(m.progressChannel))
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m installModel) View() string {
	if m.installing {
		return m.progressBar.ViewAs(m.progress)
	}
	return m.spinner.View() + " Installing..."
}

func download(plugin mtvmplugin.Plugin, version string, progressChannel chan float64) tea.Cmd {
	return func() tea.Msg {
		err := plugin.Download(version, progressChannel)
		if err != nil {
			return err
		}
		return nil
	}
}

func install(plugin mtvmplugin.Plugin, installDir string, pluginName string, version string) tea.Cmd {
	return func() tea.Msg {
		err := plugin.Install(path.Join(installDir, pluginName, version))
		if err != nil {
			return err
		}
		return installedMsg(true)
	}
}
func waitForProgress(progressChannel chan float64) tea.Cmd {
	return func() tea.Msg {
		downloadProgress := <-progressChannel
		if downloadProgress == 1 {
			return downloadedMsg(true)
		}
		return progressMsg(downloadProgress)
	}
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [tool] [version]",
	Short: "Installs a specified version of a tool",
	Long: `Installs a specified version of a tool.
For example:
If you run "mtvm install go latest" it will install the latest version of go`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var plugin mtvmplugin.Plugin
		if args[0] == "go" {
			plugin = &goplugin.Plugin{}
		} else {
			// log.Fatal bc we aren't loading plugins for now
			log.Fatal("Unknown plugin")
		}
		version := args[1]
		if version == "latest" {
			var err error
			version, err = plugin.GetLatestVersion()
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err := os.Stat(path.Join(shared.Configuration.InstallDir, args[0], version))
		if err != nil {
			if os.IsNotExist(err) {
				p := tea.NewProgram(installInitialModel(plugin, args[0], version))
				if _, err := p.Run(); err != nil {
					fmt.Printf("Alas, there's been an error: %v", err)
					os.Exit(1)
				}
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Println("That version is already installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
