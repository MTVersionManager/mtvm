package plugin

import (
	"errors"
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

// InstalledVersionCmd returns a VersionMsg on success,
// a NotFoundMsg with the plugin name if the plugin isn't found, and an error on failure
func InstalledVersionCmd(pluginName string) tea.Cmd {
	return func() tea.Msg {
		version, err := InstalledVersion(pluginName)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return err
			}
			return NotFoundMsg(pluginName)
		}
		return VersionMsg(version)
	}
}
