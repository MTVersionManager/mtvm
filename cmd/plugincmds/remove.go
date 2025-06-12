package plugincmds

import (
	"fmt"

	"github.com/MTVersionManager/mtvm/components/fatalHandler"
	"github.com/MTVersionManager/mtvm/plugin"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	StatusNone = iota
	StatusDone
	StatusNotFound
)

type removeModel struct {
	pluginName   string
	spinner      spinner.Model
	fileStatus   int
	entryStatus  int
	errorHandler fatalHandler.Model
	fileSystem afero.Fs
}

func initialRemoveModel(pluginName string) removeModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	return removeModel{
		pluginName: pluginName,
		spinner:    spin,
		fileSystem: afero.NewOsFs(),
	}
}

func (m removeModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, plugin.RemoveEntryCmd(m.pluginName, m.fileSystem), plugin.RemoveCmd(m.pluginName, m.fileSystem))
}

func (m removeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		m.errorHandler, cmd = m.errorHandler.Update(msg)
		return m, cmd
	case shared.SuccessMsg:
		if string(msg) == "RemoveEntry" {
			m.entryStatus = StatusDone
		} else if string(msg) == "Remove" {
			m.fileStatus = StatusDone
		}
	case plugin.NotFoundMsg:
		switch msg.Source {
		case "RemoveEntry":
			m.entryStatus = StatusNotFound
		case "Remove":
			m.fileStatus = StatusDone
		}
	}
	if m.entryStatus != StatusNone && m.fileStatus != StatusNone {
		return m, tea.Quit
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m removeModel) View() string {
	if m.fileStatus == StatusNone && m.entryStatus == StatusNone {
		return fmt.Sprintf("%v Removing %v...\n", m.spinner.View(), m.pluginName)
	}
	if m.fileStatus == StatusNotFound && m.entryStatus == StatusNotFound {
		return fmt.Sprintf("%v No changes were made as %v is not installed\n", shared.CheckMark, m.pluginName)
	}
	return fmt.Sprintf("%v Successfully removed %v\n", shared.CheckMark, m.pluginName)
}

var RemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a plugin",
	Long:    `Remove the plugin with the name specified`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"r", "rm"},
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialRemoveModel(args[0]))
		if model, err := p.Run(); err != nil {
			log.Fatal(err)
		} else if model, ok := model.(removeModel); ok {
			fatalHandler.Handle(model.errorHandler)
		}
	},
}
