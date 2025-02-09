package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/MTVersionManager/mtvm/shared"
	"github.com/MTVersionManager/mtvmplugin"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type removeModel struct {
	plugin     mtvmplugin.Plugin
	version    string
	pluginName string
	spinner    spinner.Model
	done       bool
}

func removeInitialModel(plugin mtvmplugin.Plugin, version, tool string) removeModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	return removeModel{
		plugin:     plugin,
		version:    version,
		pluginName: tool,
		spinner:    spin,
		done:       false,
	}
}

func (m removeModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, remove(m.plugin, m.version, filepath.Join(shared.Configuration.InstallDir, m.pluginName), shared.Configuration.PathDir))
}

func (m removeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case shared.SuccessMsg:
		if msg == "remove" {
			m.done = true
			return m, tea.Quit
		}
	case error:
		log.Fatal(msg)
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m removeModel) View() string {
	if m.done {
		return fmt.Sprintf("%v Successfully removed version %v of %v\n", shared.CheckMark, m.version, m.pluginName)
	}
	return fmt.Sprintf("%v Removing version %v of %v\n", m.spinner.View(), m.version, m.pluginName)
}

func remove(plugin mtvmplugin.Plugin, version string, installDir string, pathDir string) tea.Cmd {
	return func() tea.Msg {
		currentVer, err := plugin.GetCurrentVersion(installDir, pathDir)
		if err != nil {
			return err
		}
		err = plugin.Remove(filepath.Join(installDir, version), pathDir, version == currentVer)
		if err != nil {
			return err
		}
		return shared.SuccessMsg("remove")
	}
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [tool] [version]",
	Short: "Removes a specified version of a tool.",
	Long: `Removes a specified version of a tool.
For example:
"mtvm remove go 1.23.3" removes go version 1.23.3`,
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"r", "rm"},
	Run: func(cmd *cobra.Command, args []string) {
		plugin, err := shared.LoadPlugin(args[0])
		if err != nil {
			log.Fatal(err)
		}
		version := args[1]
		if version == "latest" {
			version, err = plugin.GetLatestVersion()
			if err != nil {
				log.Fatal(err)
			}
		}
		installed, err := shared.IsVersionInstalled(args[0], version)
		if err != nil {
			log.Fatal(err)
		}
		if installed {
			p := tea.NewProgram(removeInitialModel(plugin, version, args[0]))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("That version is not installed so you can't remove it")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
