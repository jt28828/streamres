package install

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"streamres/admin"
	"streamres/bundled"
	"streamres/file"
	"streamres/globals"
	"streamres/stdinput"
)

func Tool() error {
	// Install needs to happen as an admin
	err := CheckOrRequestAdminAccess()
	if err != nil {
		return err
	}

	// Create Cache directory first
	fmt.Println("Creating application directory:", globals.CacheDirPath)
	err = RecreateCacheDir()
	if err != nil {
		return fmt.Errorf("install failed: %s", err.Error())
	}

	fmt.Println("Moving tool to Sunshine 'tools' folder")
	err = CopyToolToSunshine()
	if err != nil {
		return fmt.Errorf("install failed: %s", err.Error())
	}

	// TODO add sunshine command prep to sunshine.conf file to save doing it manually
	// TODO add group policy addition for startup and shutdown scripts

	stdinput.AskQuestion("Install complete. Press enter or exit this window and restart sunshine to begin using streamres")
	return nil
}

func CheckOrRequestAdminAccess() error {
	if !admin.IsRunningElevatedPerms() {
		fmt.Println("Streamres install needs to run as admin to copy files to protected locations like the Streamres/tools folder.")
		stdinput.AskQuestion("You will be prompted to provide admin privileges after pressing Enter")

		// Rerun the application as an admin
		err := admin.ReRunElevatedPerms()
		if err != nil {
			return err
		}
		// Otherwise we're running again so exit the current process
		os.Exit(0)
	} else {
		fmt.Println("Installer is running as admin user, proceeding")
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

func RecreateCacheDir() error {
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
