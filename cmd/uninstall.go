package cmd

import (
	"github.com/spf13/cobra"
	"streamres/installer"
)

// uninstallCmd represents the uninstallation command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls and removes streamres configuration from sunshine",
	Long:  `Uninstalls streamres by removing all configuration directories and removing any files copied during the installation process`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return installer.Uninstall()
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
