package cmd

import (
	"log"
	"os"

	"github.com/MTVersionManager/mtvm/config"
	"github.com/MTVersionManager/mtvm/shared"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mtvm",
	Short: "A version manager that supports multiple tools.",
	Long:  `A version manager that supports multiple tools. It can do this because it has a plugins system, so you can integrate any tool you want into MTVM!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	var err error
	shared.Configuration, err = config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mtvm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
