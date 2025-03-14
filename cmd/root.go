package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"streamres/globals"
	"streamres/logging"
	"streamres/stdinput"
	"streamres/validate"
)

var verbose = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Version:           globals.Version,
	Use:               "streamres",
	Short:             "Turn on and modify the resolution of external monitors for streaming services",
	Long: `Streamres enables and disables virtual monitors and adjusts the resolution and refresh rate to suit the client device
This can be used to stream games to clients using configuration not supported by physical monitors attached to the host`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logging.Initialise(verbose)
		return validate.Application()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=====================================================================================================")
		fmt.Println("       Streamres is intended to be used from the command line. What are you intending to do?         ")
		fmt.Println("     Rerun Streamres using CMD and provide the action you want to take. eg: streamres.exe install    ")
		fmt.Println("=====================================================================================================\n\n")
		cmd.Help()
		stdinput.AskQuestion("Press enter to close")
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Whether to enable verbose logging")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
