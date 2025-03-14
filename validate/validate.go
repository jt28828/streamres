package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"streamres/bundled"
	"streamres/globals"
	"strings"
)

// Application makes sure the dependencies are initialised into the cache directory before running anything.
// This tool relies on several executables (see readme) to work so needs them places in the user cache dir
func Application() error {
	// Check if the cache directory has been initialised first
	files, err := os.ReadDir(globals.CacheDirPath)
	if err != nil || !RequiredExecutablesPresent(files) {
		return fmt.Errorf("Required installation files not found. Try running 'install' first")
	}

	current, installedVersion := installedVersionIsCurrent(files)

	if !current {
		fmt.Printf("Version mismatch found, Running version: '%s', installed version: '%s'", globals.Version, installedVersion)
		fmt.Printf("Running streamres install again is recommended to ensure up to date dependencies")
	}

	return nil
}

func RequiredExecutablesPresent(files []os.DirEntry) bool {
	for _, file := range files {
		if file.Name() == bundled.Multimonitor {
			return true
		}
	}

	return false
}

func installedVersionIsCurrent(files []os.DirEntry) (bool, string) {
	for _, file := range files {
		if file.Name() == "version" {
			contents, err := os.ReadFile(filepath.Join(globals.CacheDirPath, file.Name()))
			if err != nil {
				return false, "N/A"
			}
			version := strings.TrimSpace(string(contents))
			return version == globals.Version, version
		}
	}

	return false, "N/A"
}
