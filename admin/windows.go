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

	verbPtr := stringToUtf16Ptr(verb)
	exePtr := stringToUtf16Ptr(exe)
	argPtr := stringToUtf16Ptr(args)
	cwdPtr := stringToUtf16Ptr(cwd)

	// Execute the command in a new window
	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, windows.SW_SHOWNORMAL)
}

func stringToUtf16Ptr(input string) *uint16 {
	output, _ := syscall.UTF16PtrFromString(input)
	return output
}
