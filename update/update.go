package update

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"streamres/globals"
	"streamres/install"
	"strings"
)

func Tool() error {
	// Check if the cache directory has been initialised first
	files, err := os.ReadDir(globals.CacheDirPath)
	if err != nil || !install.RequiredExecutablesPresent(files) {
		return fmt.Errorf("Required installation files not found. Try running 'install' first")
	}
	if !installedVersionIsCurrent(files) {
		slog.Debug("Updating to new version", slog.String("version", globals.VERSION))
		install.RecreateCacheDir()
	}

	return nil
}

func installedVersionIsCurrent(files []os.DirEntry) bool {
	for _, file := range files {
		if file.Name() == "version" {
			contents, err := os.ReadFile(filepath.Join(globals.CacheDirPath, file.Name()))
			if err != nil {
				return false
			}
			return strings.TrimSpace(string(contents)) == globals.VERSION
		}
	}

	return false
}
