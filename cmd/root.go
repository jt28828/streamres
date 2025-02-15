package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"streamres/globals"
	"streamres/initialise"
)

var debugMode = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Version:           globals.VERSION,
	Use:               "streamres",
	Short:             "Turn on and modify the resolution of external monitors for streaming services",
	Long: `Streamres enables and disables virtual monitors and adjusts the resolution and refresh rate to suit the client device
This can be used to stream games to clients using configuration not supported by physical monitors attached to the host`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialise.Tool()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Whether to run in debug mode and output errors to a logfile")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
