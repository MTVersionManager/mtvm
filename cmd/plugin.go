package cmd

import (
	"fmt"

	"github.com/MTVersionManager/mtvm/plugin"

	"github.com/MTVersionManager/mtvm/cmd/plugincmds"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// pluginCmd represents the plugins command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Lists the plugins you have installed",
	Long:  `Lists the plugins you have installed`,
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := plugin.GetEntries()
		if err != nil {
			log.Fatal("Error when getting list of installed plugins", "err", err)
		}
		if len(entries) == 0 || entries == nil {
			fmt.Println("No plugins installed")
			return
		}
		for _, entry := range entries {
			fmt.Printf("%v %v\n", entry.Name, entry.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(plugincmds.InstallCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
