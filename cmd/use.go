package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [tool] [version]",
	Short: "Sets a specified version of a tool as the active version.",
	Long: `Sets a specified version of a tool as the active version.
For example:
"mtvm use go 1.23.3" sets go version 1.23.3 as the active version.
So if you run go version it will print the version number 1.23.3`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")
	useCmd.Flags().BoolP("install", "i", false, "Installs the specified version if you don't have it installed already")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
