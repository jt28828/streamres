package displays

import (
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"streamres/bundled"
	"streamres/globals"
	"strings"
	"syscall"
)

func Enable(display Display) {
	slog.Debug("Enabling display", slog.String("MonitorId", display.ShortMonitorId))
	err := runMultiMonitorToolCommand("/enable", display.ShortMonitorId)
	if err != nil {
		slog.Error("error enabling display", slog.String("MonitorId", display.ShortMonitorId), slog.String("Error", err.Error()))
	}
}

func Disable(display Display) {
	slog.Debug("Disabling display", slog.String("MonitorId", display.ShortMonitorId))
	err := runMultiMonitorToolCommand("/disable", display.ShortMonitorId)
	if err != nil {
		slog.Error("error disabling display", slog.String("MonitorId", display.ShortMonitorId), slog.String("Error", err.Error()))
	}
}

func TurnOn(display Display) {
	slog.Debug("Turning on display", slog.String("MonitorId", display.ShortMonitorId))
	err := runMultiMonitorToolCommand("/TurnOn", display.ShortMonitorId)
	if err != nil {
		slog.Error("error turning on display", slog.String("MonitorId", display.ShortMonitorId), slog.String("Error", err.Error()))
	}

}

func TurnOff(display Display) {
	slog.Debug("Turning off display", slog.String("MonitorId", display.ShortMonitorId))
	err := runMultiMonitorToolCommand("/TurnOff", display.ShortMonitorId)
	if err != nil {
		slog.Error("error turning off display", slog.String("MonitorId", display.ShortMonitorId), slog.String("Error", err.Error()))
	}
}

func SetConfig(config DisplayConfig, display Display) {
	configStr := fmt.Sprintf(`"Name=%s Width=%d Height=%d DisplayFrequency=%d"`, display.ShortMonitorId, config.Width, config.Height, config.RefreshRate)
	slog.Debug("Updating display config", slog.String("MonitorId", display.ShortMonitorId), slog.String("NewConfig", configStr))

	err := runEscapedMultiMonitorToolCommand("/SetMonitors", configStr)
	if err != nil {
		slog.Error("error updating display configuration", slog.String("MonitorId", display.ShortMonitorId), slog.String("Error", err.Error()))
	}
}

func SetAsPrimary(display Display) {
	slog.Debug("Setting display as primary", slog.String("MonitorId", display.ShortMonitorId))
	err := runEscapedMultiMonitorToolCommand("/SetPrimary", display.ShortMonitorId)
	if err != nil {
		slog.Error("error updating display configuration", slog.String("MonitorId", display.ShortMonitorId), slog.String("Error", err.Error()))
	}
}

// runEscapedMultiMonitorToolCommand runs windows commands directly using SysProcAttr to prevent invalid escaping of quotes etc.
func runEscapedMultiMonitorToolCommand(args ...string) error {
	toolPath := filepath.Join(globals.CacheDirPath, bundled.Multimonitor)
	cmd := exec.Command(toolPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.CmdLine = fmt.Sprintf(`%s %s`, toolPath, strings.Join(args, " "))
	err := cmd.Run()
	return err

}

// runEscapedMultiMonitorToolCommand runs a command via the MultiMonitorTool executable
func runMultiMonitorToolCommand(args ...string) error {
	toolPath := filepath.Join(globals.CacheDirPath, bundled.Multimonitor)
	err := exec.Command(toolPath, args...).Run()
	return err
}
