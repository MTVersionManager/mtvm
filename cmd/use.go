package cmd

import (
	"fmt"
	"log"
	"os"

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
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		installFlagUsed, err := cmd.Flags().GetBool("install")
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case len(args) == 2:
			fmt.Printf("Setting version of %v to %v\n", args[0], args[1])
		case installFlagUsed:
			fmt.Println("You need to specify a version to install.")
			err = cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(1)
		default:
			fmt.Println("I would normally list the versions available and let you pick here")
		}
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
