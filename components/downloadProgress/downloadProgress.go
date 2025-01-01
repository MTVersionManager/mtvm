package downloadProgress

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ProgressMsg   float64
	DownloadedMsg bool
)

type Model struct {
	Title           string
	progress        float64
	progressBar     progress.Model
	progressChannel chan float64
}

func New(progressChannel chan float64) Model {
	progressBar := progress.New(progress.WithDefaultGradient())
	return Model{
		progress:        0,
		progressBar:     progressBar,
		progressChannel: progressChannel,
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ProgressMsg:
		m.progress = float64(msg)
		cmd = WaitForProgress(m.progressChannel)
	case DownloadedMsg:
		m.progress = 1
	}
	return m, cmd
}

func (m Model) View() string {
	var s string
	if m.Title != "" {
		s += m.Title + "\n"
	}
	s += m.progressBar.ViewAs(m.progress)
	return s
}
