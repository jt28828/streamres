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
	"streamres/sunshine"
)

// Tool installs the streamres tool and configures sunshine to use it
func Tool() error {
	// Install needs to happen as an admin
	err := checkOrRequestAdminAccess()
	if err != nil {
		return err
	}

	// Create Cache directory first
	fmt.Println("Creating application directory:", globals.CacheDirPath)
	err = RecreateCacheDir()
	if err != nil {
		slog.Debug("Install failed to create the application directory")
		return fmt.Errorf("install failed: %s", err.Error())
	}

	// Then move this binary to the sunshine tools folder
	fmt.Println("Moving tool to Sunshine 'tools' folder")
	sunshineFolder, err := copyToolToSunshine()
	if err != nil {
		slog.Debug("Install failed to copy the streamres binary into the Sunshine 'tools' folder")
		return fmt.Errorf("install failed: %s", err.Error())
	}

	// Update sunshine config to add streamres commands to the global commands for all streaming
	err = sunshine.UpdateCommandPrep(sunshineFolder)
	if err != nil {
		slog.Debug("Install failed to update sunshines global command prep to use streamres when starting and stopping a stream")
		return fmt.Errorf("install failed: %s", err.Error())
	}

	// TODO add sunshine command prep to sunshine.conf file to save doing it manually
	// TODO add group policy addition for startup and shutdown scripts

	stdinput.AskQuestion("Install complete. Press enter or exit this window and restart sunshine to begin using streamres")
	return nil
}

func copyToolToSunshine() (string, error) {
	sunshinePath := stdinput.AskQuestionWithDefault(fmt.Sprintf("Enter your Sunshine install path. Enter to use default (%s): ", globals.DefaultSunshineInstallPath), globals.DefaultSunshineInstallPath)

	// Get local
	executable, err := os.Executable()
	if err != nil {
		slog.Debug("Failed to read running executable")
		return sunshinePath, err
	}

	// Copy file
	return sunshinePath, file.Copy(executable, filepath.Join(sunshinePath, "tools", globals.StreamresExecutableName))
}

func checkOrRequestAdminAccess() error {
	if !admin.IsRunningElevatedPerms() {
		fmt.Println("Streamres install needs to run as admin to copy files to protected locations like the Streamres/tools folder.")
		stdinput.AskQuestion("You will be prompted to provide admin privileges after pressing Enter")

		// Rerun the application as an admin
		err := admin.ReRunWithElevatedPerms()
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

func RecreateCacheDir() error {
	// Remove just in case it's an update
	_ = os.Remove(globals.CacheDirPath)

	// Create folder and get a reference
	err := os.MkdirAll(globals.CacheDirPath, 0755)
	if err != nil {
		return err
	}

	// Add required dependencies to the cache folder
	err = os.WriteFile(filepath.Join(globals.CacheDirPath, "version"), []byte(globals.Version), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(globals.CacheDirPath, bundled.Multimonitor), bundled.LoadMultiMonitorTool(), 0755)
	if err != nil {
		return err
	}

	return nil
}
