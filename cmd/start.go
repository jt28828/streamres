package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"streamres/displays"
	"strings"
	"time"
)

var (
	width       int
	height      int
	refreshRate int
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Enables the virtual monitor and sets configuration",
	Long: `start enables the virtual monitor and sets the resolution and refresh rate requested by the client.
Disables any attached physical monitors while in use and stores the state of the system before starting to allow a revert`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get monitor state first
		state := displays.GetCurrentState()
		if state == nil {
			return fmt.Errorf("cannot retrieve monitor state, something went wrong calling MultiMonitorTool.exe")
		}

		// Find the virtual monitor
		var virtualDisplay *displays.Display = nil
		var otherDisplays []displays.Display
		for _, display := range state {
			if strings.Contains(display.Adapter, "Virtual Display Driver") {
				virtualDisplay = &display
			} else {
				otherDisplays = append(otherDisplays, display)
			}
		}
		if virtualDisplay == nil {
			return fmt.Errorf("cannot find Virtual Display Driver. Make sure you have installed https://github.com/itsmikethetech/Virtual-Display-Driver")
		}

		// Found virtual monitor. Enable, set the resolution and mark it as the primary monitor.
		displays.Enable(*virtualDisplay)
		displays.SetConfig(displays.DisplayConfig{Width: width, Height: height, RefreshRate: refreshRate}, *virtualDisplay)
		// Wait for changes to settle otherwise monitors can end up flickering
		time.Sleep(1 * time.Second)
		displays.SetAsPrimary(*virtualDisplay)
		time.Sleep(1 * time.Second)

		// Finally turn off all other monitors
		for _, display := range otherDisplays {
			displays.TurnOff(display)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVar(&width, "width", 1920, "Target width for the virtual display")
	startCmd.Flags().IntVar(&height, "height", 1080, "Target height for the virtual display")
	startCmd.Flags().IntVar(&refreshRate, "refresh", 60, "Target refresh rate for the virtual display")
}
