package plugincmds

import (
	"encoding/json"
	"fmt"
	"github.com/MTVersionManager/mtvm/components/downloader"
	"github.com/MTVersionManager/mtvm/shared"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type installModel struct {
	downloader downloader.Model
	step       int
}

func initialInstallModel(url string) installModel {
	return installModel{
		downloader: downloader.New(url, downloader.UseTitle("Downloading plugin metadata...")),
	}
}

func loadMetadata(rawData []byte) (shared.PluginMetadata, error) {
	var metadata shared.PluginMetadata
	err := json.Unmarshal(rawData, &metadata)
	if err != nil {
		return metadata, err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}

func loadMetadataCmd(rawData []byte) tea.Cmd {
	return func() tea.Msg {
		metadata, err := loadMetadata(rawData)
		if err != nil {
			return err
		}
		return metadata
	}
}

func (m installModel) Init() tea.Cmd {
	return m.downloader.Init()
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case shared.SuccessMsg:
		if msg == "download" {
			log.Printf("Downloaded data length: %d bytes (install)", len(msg))
			cmds = append(cmds, loadMetadataCmd(m.downloader.GetDownloadedData()))
		}
	case shared.PluginMetadata:
		fmt.Println(msg)
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m installModel) View() string {
	return m.downloader.View() + "\n"
}

var InstallCmd = &cobra.Command{
	Use:   "install [plugin url]",
	Short: "Install a plugin",
	Long:  `Install a plugin given a link to the plugin's metadata JSON`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Var(args[0], "http_url")
		if err != nil {
			fmt.Println("Please enter a valid http url")
			os.Exit(1)
		}
		p := tea.NewProgram(initialInstallModel(args[0]))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
		}
	},
}
