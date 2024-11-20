package cmd

import (
	"fmt"
	"github.com/MTVersionManager/mtvm/components/download"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/MTVersionManager/mtvmplugin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type installModel struct {
	downloader download.Model
}

func installInitialModel(plugin mtvmplugin.Plugin, pluginName string, version string) installModel {
	downloadModel := download.New(plugin, pluginName, version)
	return installModel{
		downloader: downloadModel,
	}
}

func (m installModel) Init() tea.Cmd {
	return m.downloader.Init()
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.downloader, cmd = m.downloader.Update(msg)
	return m, cmd
}

func (m installModel) View() string {
	return m.downloader.View()
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
		plugin, err := shared.LoadPlugin(args[0])
		if err != nil {
			log.Fatal(err)
		}
		version := args[1]
		if version == "latest" {
			var err error
			version, err = plugin.GetLatestVersion()
			if err != nil {
				log.Fatal(err)
			}
		}
		installed, err := shared.IsVersionInstalled(args[0], version)
		if err != nil {
			log.Fatal(err)
		}
		if !installed {
			p := tea.NewProgram(installInitialModel(plugin, args[0], version))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("That version is already installed")
			os.Exit(1)
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
