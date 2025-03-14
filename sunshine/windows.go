//go:build windows

package sunshine

import (
	"path/filepath"
	"strings"
)

func addStreamresCommands(commands []PrepCommand, sunshineFolder string) []PrepCommand {
	streamresInstallPath := filepath.Join(sunshineFolder, "tools", "streamres.exe")

	// Command to setup and tear down virtual monitor
	startAndStopCommands := PrepCommand{
		Do:       `cmd /C ` + streamresInstallPath + ` start --height %SUNSHINE_CLIENT_HEIGHT% --width %SUNSHINE_CLIENT_WIDTH% --refresh %SUNSHINE_CLIENT_FPS%`,
		Undo:     `cmd /C ` + streamresInstallPath + ` revert`,
		Elevated: false,
	}

	// Command to wait after monitor configuration to prevent applications launching on the wrong monitor while it's turning off
	waitCommand := PrepCommand{
		Do:       "timeout /t 3",
		Undo:     "",
		Elevated: false,
	}

	commands = append(commands, startAndStopCommands, waitCommand)
	return commands
}

func removeStreamresCommands(commands []PrepCommand) []PrepCommand {
	newCommands := []PrepCommand{}
	for _, command := range commands {
		if !strings.Contains(command.Do, "streamres.exe") && command.Do != "timeout /t 3" {
			newCommands = append(newCommands, command)
		}
	}

	return newCommands
}
