package main

import (
	"github.com/spf13/cobra"
	"streamres/cmd"
)

func main() {
	// Turn off mousetrap for windows
	cobra.MousetrapHelpText = ""

	// Start cli
	cmd.Execute()
}
