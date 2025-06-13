package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MTVersionManager/mtvm/components/install"
	"github.com/MTVersionManager/mtvm/shared"
	"github.com/MTVersionManager/mtvmplugin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type installModel struct {
	installer  install.Model
	installed  bool
	version    string
	pluginName string
}

func installInitialModel(plugin mtvmplugin.Plugin, pluginName, version string) installModel {
	downloadModel := install.New(plugin, pluginName, version)
	return installModel{
		installer:  downloadModel,
		installed:  false,
		version:    version,
		pluginName: pluginName,
	}
}

func (m installModel) Init() tea.Cmd {
	return m.installer.Init()
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg.(install.InstalledMsg) {
		m.installed = true
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.installer, cmd = m.installer.Update(msg)
	return m, cmd
}

func (m installModel) View() string {
	if !m.installed {
		return m.installer.View()
	}
	return fmt.Sprintf("%v Installed version %v of %v\n", shared.CheckMark, m.version, m.pluginName)
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [tool] [version]",
	Short: "Installs a specified version of a tool",
	Long: `Installs a specified version of a tool.
For example:
If you run "mtvm install go latest" it will install the latest version of go`,
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"i", "in"},
	Run: func(cmd *cobra.Command, args []string) {
		err := createInstallDir(afero.NewOsFs())
		if err != nil {
			log.Fatal(err)
		}
		plugin, err := shared.LoadPlugin(args[0])
		if err != nil {
			log.Fatal(err)
		}
		version := args[1]
		if strings.ToLower(version) == "latest" {
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
				log.Fatal(err)
			}
		} else {
			fmt.Println("That version is already installed")
			os.Exit(1)
		}
	},
}

func createInstallDir(fs afero.Fs) error {
	err := fs.MkdirAll(shared.Configuration.InstallDir, 0o777)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
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
