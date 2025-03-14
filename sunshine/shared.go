package sunshine

import (
	"bufio"
	"bytes"
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
	sunshineConfFilepath := filepath.Join(sunshineFolder, "config", "sunshine.conf")

	prepCommands, restOfFile, err := readCommandPrep(sunshineConfFilepath)

	replacementCommands, err := getUpdatedCommandsEntry(prepCommands, sunshineFolder)

	if err != nil {
		slog.Debug("Failed to update existing global_prep_cmd entry")
		return err
	}

	// Convert the replacement commands back to a single line string
	commandsBytes := &bytes.Buffer{}
	encoder := json.NewEncoder(commandsBytes)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(replacementCommands)

	if err != nil {
		slog.Debug("Failed to marshal commands to string")
		return err
	}

	// The JSON encoder adds a newline itself after the contents. So directly append the rest of the file
	output := fmt.Sprintf("global_prep_cmd = %s", commandsBytes.String()) + restOfFile

	// Overwrite the file
	err = os.WriteFile(sunshineConfFilepath, []byte(output), 0644)
	if err != nil {
		slog.Debug("Failed to write config back into file")
		return err
	}

	return nil
}

func getUpdatedCommandsEntry(prepCommands, sunshineFolder string) ([]PrepCommand, error) {
	// If there are no commands we just init an empty array
	if len(prepCommands) == 0 {
		prepCommands = "[]"
	}

	// Convert so we can add our command
	commands := []PrepCommand{}

	err := json.Unmarshal([]byte(prepCommands), &commands)
	if err != nil {
		return nil, err
	}

	// Strip out any streamres commands if they are already in there to prevent duplicates
	commands = removeStreamresCommands(commands)

	// Now add the new streamres commands to the array
	commands = addStreamresCommands(commands, sunshineFolder)

	return commands, nil
}

func readCommandPrep(filepath string) (string, string, error) {
	// Read each line of the config file and extract the command prep
	sunshineConfFile, err := os.Open(filepath)

	if err != nil {
		slog.Debug("Failed to read sunshine config file")
		return "", "", err
	}

	prepCommands := ""
	restOfFile := &strings.Builder{}
	scanner := bufio.NewScanner(sunshineConfFile)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "global_prep_cmd") {
			prepCommands = strings.TrimPrefix(scanner.Text(), "global_prep_cmd =")
		} else {
			restOfFile.WriteString(scanner.Text() + "\n")
		}
	}

	// Strip leading and trailing newlines
	trimmedFile := strings.Trim(restOfFile.String(), "\n")
	return prepCommands, trimmedFile, nil
}
