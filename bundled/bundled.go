package bundled

import (
	"embed"
)

const Multimonitor = "MultiMonitorTool.exe"

//go:embed executable/*
var executables embed.FS

func LoadMultiMonitorTool() []byte {
	file, err := executables.ReadFile("executable/MultiMonitorTool.exe")
	if err != nil {
		panic(err)
	}
	return file
}
