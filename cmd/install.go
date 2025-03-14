package cmd

import (
	"github.com/spf13/cobra"
	"streamres/install"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs and configures sunshine to use streamres",
	Long: `Installs streamres by creating configuration directories and moving some files. 
This includes installing dependencies and setting up a cache folder, 
copying streamres to the sunshine tools folder, 
and setting up system startup and shutdown scripts to allow clean recovery of monitor configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return install.Tool()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
