package initialise

import (
	"os"
	"path/filepath"
	"streamres/bundled"
	"streamres/globals"
	"strings"
)

// Tool makes sure the dependencies are initialised into the cache directory before running anything.
// This tool relies on several executables (see readme) to work so needs them places in the user cache dir
func Tool() error {
	// Check if the cache directory has been initialised
	// Create Application Cache Dir if missing
	files, err := os.ReadDir(globals.CacheDirPath)

	if err == nil {
		// Folder exists. Make sure what we need is in it
		if installedVersionIsCurrent(files) && requiredExecutablesPresent(files) {
			// All good to continue
			return nil
		}
	}

	// Otherwise we need to recreate the directory because it's missing or out of date
	return recreateCacheDir()
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

func requiredExecutablesPresent(files []os.DirEntry) bool {
	for _, file := range files {
		if file.Name() == bundled.Multimonitor {
			return true
		}
	}

	return false
}

func recreateCacheDir() error {
	// Remove just in case it's an update
	_ = os.Remove(globals.CacheDirPath)

	// Create folder and get a reference
	err := os.MkdirAll(globals.CacheDirPath, 0755)
	if err != nil {
		return err
	}

	// Add required dependencies to the cache folder
	err = os.WriteFile(filepath.Join(globals.CacheDirPath, "version"), []byte(globals.VERSION), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(globals.CacheDirPath, bundled.Multimonitor), bundled.LoadMultiMonitorTool(), 0755)
	if err != nil {
		return err
	}

	return nil
}
