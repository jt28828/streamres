//go:build windows

package admin

import (
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
)

func IsRunningElevatedPerms() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")

	return err == nil
}

// ReRunWithElevatedPerms runs the application again as an admin
// See https://learn.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shellexecutea
func ReRunWithElevatedPerms() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	// Show the command in a new window
	var showCmd int32 = windows.SW_SHOWNORMAL

	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
}
