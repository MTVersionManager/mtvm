package plugin

import (
	"github.com/MTVersionManager/mtvm/shared"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
)

// UpdateEntriesCmd returns an error on failure and a shared.SuccessMsg with contents "UpdateEntries" on success
func UpdateEntriesCmd(entry Entry, fs afero.Fs) tea.Cmd {
	return func() tea.Msg {
		err := UpdateEntries(entry, fs)
		if err != nil {
			return err
		}
		return shared.SuccessMsg("UpdateEntries")
	}
}

// InstalledVersionCmd returns a VersionMsg on success,
// a NotFoundMsg with the plugin name if the plugin isn't found, and an error on failure
func InstalledVersionCmd(pluginName string, fs afero.Fs) tea.Cmd {
	return func() tea.Msg {
		version, err := InstalledVersion(pluginName, fs)
		if err != nil {
			if shared.IsNotFound(err) {
				return NotFoundMsg{
					PluginName: pluginName,
					Source:     "InstalledVersion",
				}
			}
			return err
		}
		return VersionMsg(version)
	}
}

func RemoveEntryCmd(pluginName string, fs afero.Fs) tea.Cmd {
	return func() tea.Msg {
		err := RemoveEntry(pluginName, fs)
		if err != nil {
			if shared.IsNotFound(err) {
				return NotFoundMsg{
					PluginName: pluginName,
					Source:     "RemoveEntry",
				}
			}
			return err
		}
		return shared.SuccessMsg("RemoveEntry")
	}
}

func RemoveCmd(pluginName string, fs afero.Fs) tea.Cmd {
	return func() tea.Msg {
		err := Remove(pluginName, fs)
		if err != nil {
			if shared.IsNotFound(err) {
				return NotFoundMsg{
					PluginName: pluginName,
					Source:     "Remove",
				}
			}
			return err
		}
		return shared.SuccessMsg("Remove")
	}
}
