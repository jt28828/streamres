package sunshine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type PrepCommand struct {
	Do       string `json:"do"`
	Undo     string `json:"undo"`
	Elevated bool   `json:"elevated"`
}

// UpdateCommandPrep updates the sunshine global_prep_cmd setting to add streamres entries
func UpdateCommandPrep(sunshineFolder string) error {
	// Get file contents first
	sunshineConfFile, err := os.Open(filepath.Join(sunshineFolder, "config", "sunshine.conf"))
	defer sunshineConfFile.Close()

	if err != nil {
		slog.Debug("Failed to read sunshine config file")
		return err
	}

	prepCommands, restOfFile := readCommandPrep(sunshineConfFile)

	// Convert so we can add our command
	commands := []PrepCommand{}

	err = json.Unmarshal([]byte(prepCommands), &commands)
	if err != nil {
		slog.Debug("Failed to parse existing global_prep_cmd entry")
		return err
	}

	// Strip out any streamres commands if they are already in there
	commands = removeStreamresCommands(commands)
	// Now add the new streamres commands to the array
	commands = addStreamresCommands(commands)

	// Convert the entry back to a single line string
	commandsStr, err := json.Marshal(commands)
	if err != nil {
		slog.Debug("Failed to marshal commands to string")
		return err
	}

	output := fmt.Sprintf("global_prep_cmd = %s", string(commandsStr)) + "\n" + restOfFile

	// Overwrite the file
	_, err = sunshineConfFile.Seek(0, 0)
	if err != nil {
		slog.Debug("Failed to seek back to the start of the file")
		return err
	}

	_, err = sunshineConfFile.Write([]byte(output))
	if err != nil {
		slog.Debug("Failed to write config back into file")
		return err
	}

	return nil
}

func addStreamresCommands(commands []PrepCommand) []PrepCommand {
	startAndStopCommands := PrepCommand{
		Do:       `cmd /C "C:\Program Files\Sunshine\tools\streamres.exe" start --height %SUNSHINE_CLIENT_HEIGHT% --width %SUNSHINE_CLIENT_WIDTH% --refresh %SUNSHINE_CLIENT_FPS%`,
		Undo:     `cmd /C "C:\Program Files\Sunshine\tools\streamres.exe" revert`,
		Elevated: false,
	}
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
		if !strings.Contains(command.Do, "streamres") && command.Do != "timeout /t 3" {
			newCommands = append(newCommands, command)
		}
	}

	return newCommands
}

func readCommandPrep(configFile *os.File) (prepCommands, restOfFile string) {
	// Read each line of the config file and extract the command prep
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "global_prep_cmd") {
			prepCommands = strings.TrimPrefix(scanner.Text(), "global_prep_cmd =")
		} else {
			restOfFile += scanner.Text() + "\n"
		}
	}
	return
}
