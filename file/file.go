package file

import (
	"log/slog"
	"os"
)

func Copy(source string, destination string) error {
	sourceContents, err := os.ReadFile(source)
	if err != nil {
		slog.Debug("Failed to read file contents", slog.String("path", source))
		return err
	}

	destinationFile, err := os.Create(destination)
	defer destinationFile.Close()

	if err != nil {
		slog.Debug("Failed to open destination file")
		return err
	}

	// Always overwrite the destination file
	_, err = destinationFile.Write(sourceContents)

	if err != nil {
		slog.Debug("Failed to write data to destination file")
		return err
	}

	return nil
}
