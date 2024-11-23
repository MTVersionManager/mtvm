package install

import (
	"github.com/MTVersionManager/mtvm/components/downloadProgress"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/MTVersionManager/mtvmplugin"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"path/filepath"
)

type Model struct {
	progressChannel chan float64
	plugin          mtvmplugin.Plugin
	installing      bool
	version         string
	spinner         spinner.Model
	pluginName      string
	downloader      downloadProgress.Model
}

type InstalledMsg bool

func New(plugin mtvmplugin.Plugin, pluginName string, version string) Model {
	progressChannel := make(chan float64)
	downloader := downloadProgress.New(progressChannel)
	downloader.Title = "Downloading..."
	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	return Model{
		progressChannel: progressChannel,
		version:         version,
		plugin:          plugin,
		downloader:      downloader,
		installing:      true,
		spinner:         spinnerModel,
		pluginName:      pluginName,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(downloadProgress.WaitForProgress(m.progressChannel), Download(m.plugin, m.version, m.progressChannel), m.spinner.Tick)
}

func Download(plugin mtvmplugin.Plugin, version string, progressChannel chan float64) tea.Cmd {
	return func() tea.Msg {
		err := plugin.Download(version, progressChannel)
		if err != nil {
			return err
		}
		return nil
	}
}

func Install(plugin mtvmplugin.Plugin, installDir string, pluginName string, version string) tea.Cmd {
	return func() tea.Msg {
		err := plugin.Install(filepath.Join(installDir, pluginName, version))
		if err != nil {
			return err
		}
		return InstalledMsg(true)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case error:
		log.Fatal(msg)
	case downloadProgress.DownloadedMsg:
		m.installing = false
		cmds = append(cmds, Install(m.plugin, shared.Configuration.InstallDir, m.pluginName, m.version))
	case InstalledMsg:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.installing {
		return m.downloader.View()
	}
	return m.spinner.View() + " Installing..."
}
