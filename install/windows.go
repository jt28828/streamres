//go:build windows

package install

import (
	"log/slog"
	"os"
	"path/filepath"
	"streamres/file"
	"streamres/stdinput"
)

func CopyToolToSunshine() error {
	sunshinePath := stdinput.AskQuestionWithDefault(`Enter your Sunshine install path. Enter to use default (C:\Program Files\Sunshine): `, `C:\Program Files\Sunshine`)

	// Get local
	executable, err := os.Executable()
	if err != nil {
		slog.Debug("Failed to read running executable")
		return err
	}

	// Copy file
	return file.Copy(executable, filepath.Join(sunshinePath, "tools", "streamres.exe"))
}
