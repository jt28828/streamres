package displays

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"streamres/bundled"
	"streamres/globals"
	"strings"
	"syscall"
)

func Enable(display Display) {
	err := runMultiMonitorToolCommand("/enable", display.ShortMonitorId)
	if err != nil {
		log.Println(fmt.Errorf("error enabling display %s: %s", display.ShortMonitorId, err.Error()))
	}
}

func Disable(display Display) {
	err := runMultiMonitorToolCommand("/disable", display.ShortMonitorId)
	if err != nil {
		log.Println(fmt.Errorf("error disabling display %s: %s", display.ShortMonitorId, err.Error()))
	}
}

func TurnOn(display Display) {
	err := runMultiMonitorToolCommand("/TurnOn", display.ShortMonitorId)
	if err != nil {
		log.Println(fmt.Errorf("error turning on display %s: %s", display.ShortMonitorId, err.Error()))
	}

}

func TurnOff(display Display) {
	err := runMultiMonitorToolCommand("/TurnOff", display.ShortMonitorId)
	if err != nil {
		log.Println(fmt.Errorf("error turning off display %s: %s", display.ShortMonitorId, err.Error()))
	}

}

func SetConfig(config DisplayConfig, display Display) {
	configStr := fmt.Sprintf(`"Name=%s Width=%d Height=%d DisplayFrequency=%d"`, display.ShortMonitorId, config.Width, config.Height, config.RefreshRate)
	err := runEscapedMultiMonitorToolCommand("/SetMonitors", configStr)
	if err != nil {
		log.Println(fmt.Errorf("error updating display configuration %s: %s", display.ShortMonitorId, err.Error()))
	}
}

func SetAsPrimary(display Display) {
	err := runEscapedMultiMonitorToolCommand("/SetPrimary", display.ShortMonitorId)
	if err != nil {
		log.Println(fmt.Errorf("error updating display configuration %s: %s", display.ShortMonitorId, err.Error()))
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
