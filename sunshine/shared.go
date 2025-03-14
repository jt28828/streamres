package sunshine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"streamres/file"
	"streamres/globals"
	"streamres/stdinput"
	"strings"
)

type PrepCommand struct {
	Do       string `json:"do"`
	Undo     string `json:"undo"`
	Elevated bool   `json:"elevated"`
}

func Install() error {
	sunshineFolder, streamresExePath := getSunshineInstallPaths()

	// First Copy streamres to sunshine
	err := copyToolToSunshine(streamresExePath)
	if err != nil {
		slog.Debug("Failed to copy streamres.exe to sunshine 'tools' folder")
		return err
	}
	fmt.Println("[X] - Copied streamres to Sunshine tools folder", streamresExePath)

	// Then update sunshine conf to use streamres for every stream
	err = addOrRemoveCommandPrep(sunshineFolder, true)
	if err != nil {
		slog.Debug("Failed to update the global_prep_cmd entry in sunshine conf")
	}
	fmt.Println("[X] - Updated sunshine.conf file to add streamres to the global_prep_cmd list")

	return err
}

func Uninstall() error {
	sunshineFolder, streamresExePath := getSunshineInstallPaths()

	// First delete streamres from the sunshine 'tools' folder
	err := removeStreamresTool(streamresExePath)
	if err != nil {
		slog.Debug("Failed to delete streamres.exe from sunshine 'tools' folder")
		return err
	}

	// Then update sunshine conf to remove streamres from the global prep commands
	err = addOrRemoveCommandPrep(sunshineFolder, false)
	if err != nil {
		slog.Debug("Failed to update the global_prep_cmd entry in sunshine conf")
	}

	return err
}

// getSunshineInstallPaths retrieves the sunshine root installation folder, and the target filepath for the streamres executable within the tools folder
func getSunshineInstallPaths() (sunshineFolder, targetExecutablePath string) {
	sunshineFolder = stdinput.AskQuestionWithDefault(fmt.Sprintf("Enter your Sunshine install path. Enter to use default (%s): ", globals.DefaultSunshineInstallPath), globals.DefaultSunshineInstallPath)
	fmt.Println("\rUsing", sunshineFolder, "as the Sunshine install path")
	targetExecutablePath = filepath.Join(sunshineFolder, "tools", globals.StreamresExecutableName)
	return
}

// removeStreamresTool removes streamres.exe from the sunshine `tools` folder if present. If not already present the action is a no-op
func removeStreamresTool(streamresExePath string) error {
	if _, err := os.Stat(streamresExePath); errors.Is(err, os.ErrNotExist) {
		// Already deleted or never installed
		return nil
	}

	// Found the file so we need to remove it
	return os.Remove(streamresExePath)
}

func copyToolToSunshine(targetExecutablePath string) error {
	// Get local
	executable, err := os.Executable()
	if err != nil {
		slog.Debug("Failed to read running executable")
		return err
	}

	// Copy file
	return file.Copy(executable, targetExecutablePath)
}

// addOrRemoveCommandPrep updates the sunshine global_prep_cmd setting to add or remove streamres entries
func addOrRemoveCommandPrep(sunshineFolder string, addCommands bool) error {
	sunshineConfFilepath := filepath.Join(sunshineFolder, "config", "sunshine.conf")

	prepCommands, restOfFile, err := readCommandPrep(sunshineConfFilepath)

	replacementCommands, err := getUpdatedCommandsEntry(prepCommands, sunshineFolder, addCommands)

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

func getUpdatedCommandsEntry(prepCommands, sunshineFolder string, addCommands bool) ([]PrepCommand, error) {
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

	if addCommands {
		// Now add the new streamres commands to the array
		commands = addStreamresCommands(commands, sunshineFolder)
	}

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
