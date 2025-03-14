package installer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"streamres/admin"
	"streamres/bundled"
	"streamres/globals"
	"streamres/stdinput"
	"streamres/sunshine"
)

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

// Install installs the streamres tool and configures sunshine to use it
func Install() error {
	// Install needs to happen as an admin
	err := checkOrRequestAdminAccess()
	if err != nil {
		return err
	}

	// Create Cache directory first
	err = RecreateCacheDir()
	if err != nil {
		slog.Debug("Failed to gain admin access to continue install")
		return fmt.Errorf("install failed: %s", err.Error())
	}
	fmt.Println("[X] - Created application directory", globals.CacheDirPath)

	// Then run the sunshine folder installation tasks
	err = sunshine.Install()
	if err != nil {
		slog.Debug("Install failed to run actions in the Sunshine installation folder")
		return fmt.Errorf("install failed: %s", err.Error())
	}

	// TODO maybe add group policy addition for startup and shutdown scripts to revert streaming if PC is shut down from sunshine stream

	stdinput.AskQuestion("\nInstall complete. Press enter or exit this window and restart sunshine to begin using streamres")
	return nil
}

// Uninstall uninstalls the streamres tool
func Uninstall() error {
	// Uninstall needs to happen as an admin
	err := checkOrRequestAdminAccess()
	if err != nil {
		slog.Debug("Failed to gain admin access to continue uninstall")
		return err
	}

	// Destroy Cache directory first
	err = os.RemoveAll(globals.CacheDirPath)
	if err != nil {
		slog.Debug("Failed to remove cache directory during uninstall")
		return err
	}
	fmt.Println("[X] - Deleted application directory", globals.CacheDirPath)

	// Then run the sunshine folder uninstallation tasks
	err = sunshine.Uninstall()
	if err != nil {
		slog.Debug("Uninstall failed to run actions in the Sunshine installation folder, streamres may still be partially installed")
		return fmt.Errorf("install failed: %s", err.Error())
	}

	stdinput.AskQuestion("\nStreamres removed successfully. Press enter or exit this window")
	return nil
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
