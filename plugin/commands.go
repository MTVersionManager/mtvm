package plugin

import (
	"github.com/MTVersionManager/mtvm/shared"
	tea "github.com/charmbracelet/bubbletea"
)

// UpdateEntriesCmd returns an error on failure and a shared.SuccessMsg with contents "UpdateEntries" on success
func UpdateEntriesCmd(entry Entry) tea.Cmd {
	return func() tea.Msg {
		err := UpdateEntries(entry)
		if err != nil {
			return err
		}
		return shared.SuccessMsg("UpdateEntries")
	}
}
