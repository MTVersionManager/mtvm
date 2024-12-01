package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// PluginsCmd represents the plugins command
var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Lists the plugins you have installed",
	Long:  `Lists the plugins you have installed`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I would be listing the plugins you have installed")
	},
}

func init() {
	rootCmd.AddCommand(PluginsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
