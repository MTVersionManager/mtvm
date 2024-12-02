package plugincmds

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"os"
)

var InstallCmd = &cobra.Command{
	Use:   "install [plugin url]",
	Short: "Install a plugin",
	Long:  `Install a plugin given the link to the plugin's metadata JSON'`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Install plugin called")
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Var(args[0], "http_url")
		if err != nil {
			fmt.Println("Please enter a valid http url")
			os.Exit(1)
		}
	},
}
