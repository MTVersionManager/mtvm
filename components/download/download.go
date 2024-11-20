package download

import (
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/MTVersionManager/mtvmplugin"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"path"
)

type Model struct {
	progressChannel chan float64
	plugin          mtvmplugin.Plugin
	installing      bool
	progress        float64
	progressBar     progress.Model
	version         string
	spinner         spinner.Model
	pluginName      string
}

type ProgressMsg float64
type DownloadedMsg bool
type InstalledMsg bool

func New(plugin mtvmplugin.Plugin, pluginName string, version string) Model {
	progressChannel := make(chan float64)
	progressBar := progress.New(progress.WithDefaultGradient())
	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	return Model{
		progressChannel: progressChannel,
		progressBar:     progressBar,
		version:         version,
		plugin:          plugin,
		installing:      true,
		spinner:         spinnerModel,
		pluginName:      pluginName,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(WaitForProgress(m.progressChannel), Download(m.plugin, m.version, m.progressChannel), m.spinner.Tick)
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
		err := plugin.Install(path.Join(installDir, pluginName, version))
		if err != nil {
			return err
		}
		return InstalledMsg(true)
	}
}

func WaitForProgress(progressChannel chan float64) tea.Cmd {
	return func() tea.Msg {
		downloadProgress := <-progressChannel
		if downloadProgress == 1 {
			return DownloadedMsg(true)
		}
		return ProgressMsg(downloadProgress)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case error:
		log.Fatal(msg)
	case DownloadedMsg:
		m.progress = 1
		m.installing = false
		return m, Install(m.plugin, shared.Configuration.InstallDir, m.pluginName, m.version)
	case InstalledMsg:
		return m, tea.Quit
	case ProgressMsg:
		m.progress = float64(msg)
		cmds = append(cmds, WaitForProgress(m.progressChannel))
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.installing {
		return "Downloading...\n" + m.progressBar.ViewAs(m.progress)
	}
	return m.spinner.View() + " Installing..."
}
