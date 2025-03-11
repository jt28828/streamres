package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"streamres/displays"
	"strings"
)

// revertCmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Reverts display settings to what they were before the previous 'start' command was run",
	Long: `Reverts cached display settings to what they were previously
before the virtual monitor was enabled and configured.
Restores usage of hardware monitors while disabling the virtual display`,
	RunE: func(cmd *cobra.Command, args []string) error {
		state := displays.GetPreviousState()
		if state == nil {
			return fmt.Errorf("no previous state now, start command may not have run yet")
		}

		// Restore the previous primary display first
		var previousPrimaryDisplay *displays.Display
		var virtualDisplay *displays.Display
		var otherDisplays []displays.Display
		for _, display := range state {
			slog.Debug("Found display", slog.String("name", display.Adapter), slog.String("id", display.ShortMonitorId))
			if display.PrimaryMonitor {
				previousPrimaryDisplay = &display
			} else if strings.Contains(display.Adapter, "Virtual Display Driver") {
				virtualDisplay = &display
			} else {
				otherDisplays = append(otherDisplays, display)
			}
		}

		if previousPrimaryDisplay == nil {
			log.Println("no previous primary display found in state, falling back to first monitor in list that isn't the virtual display adapter")
			for _, display := range state {
				if !strings.Contains(display.Adapter, "Virtual Display Driver") {
					previousPrimaryDisplay = &display
					break
				}
			}
		}

		// Set the primary display back to what it was
		displays.TurnOn(*previousPrimaryDisplay)
		displays.SetAsPrimary(*previousPrimaryDisplay)
		displays.SetConfig(previousPrimaryDisplay.Config, *previousPrimaryDisplay)

		// Disable the virtual monitor
		displays.Disable(*virtualDisplay)

		// Restore the state of the other displays
		for _, display := range otherDisplays {
			if !display.Disconnected {
				displays.TurnOn(display)
			}
		}

		// Finally purge the state file to prevent reverting again
		return displays.DeleteStateFile()
	},
}

func init() {
	rootCmd.AddCommand(revertCmd)
}
