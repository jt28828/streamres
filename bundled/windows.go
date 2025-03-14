//go:build windows

package bundled

import (
	"embed"
)

const Multimonitor = "MultiMonitorTool.exe"

//go:embed windowsExecutable/*
var executables embed.FS

func LoadMultiMonitorTool() []byte {
	file, err := executables.ReadFile("windowsExecutable/MultiMonitorTool.exe")
	if err != nil {
		panic(err)
	}
	return file
}
