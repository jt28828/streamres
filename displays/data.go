package displays

import (
	"bytes"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"streamres/globals"
	"strings"
)

type DisplayConfig struct {
	Width       int `json:"width"`
	Height      int `json:"height"`
	RefreshRate int `json:"refreshRate"`
}
type Display struct {
	ShortMonitorId string        `json:"shortMonitorId"`
	Adapter        string        `json:"adapter"`
	Active         bool          `json:"active"`
	Disconnected   bool          `json:"disconnected"`
	PrimaryMonitor bool          `json:"primaryMonitor"`
	Config         DisplayConfig `json:"config"`
}

func GetCurrentState() []Display {
	state := getMonitorInfoCsv()
	if state == nil {
		return nil
	}

	return formatCsvToStructs(state)
}

func GetPreviousState() []Display {
	state := loadMonitorTsvState()
	if state == nil {
		return nil
	}

	return formatCsvToStructs(state)
}

func getMonitorInfoCsv() []byte {
	// Export monitor info using MultiMonitorTool
	tsvOutputPath := filepath.Join(globals.CacheDirPath, "monitors.tsv")
	if err := runMultiMonitorToolCommand("/stab", tsvOutputPath); err != nil {
		return nil
	}

	// Command ran successfully. Return the file contents
	return loadMonitorTsvState()
}

func loadMonitorTsvState() []byte {
	tsvOutputPath := filepath.Join(globals.CacheDirPath, "monitors.tsv")
	file, _ := os.ReadFile(tsvOutputPath)

	// File is in UTF16, need to convert to UTF8 to work with it

	win16Enc := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16Bom := unicode.BOMOverride(win16Enc.NewDecoder())

	// Make a reader that can read utf16
	reader := transform.NewReader(bytes.NewReader(file), utf16Bom)

	decoded, err := io.ReadAll(reader)

	if err != nil {
		log.Println("Error reading from TSV state file:", err)
	}

	return decoded
}

// formatCsvToStructs takes a CSV file output by the MultiMonitorTool and converts it in go structs for easy access
func formatCsvToStructs(stateCsv []byte) []Display {
	displays := []Display{}
	// Split by line to get each individual monitor

	displaysLine := strings.Split(string(stateCsv), "\n")

	// Skip the first and last rows because they are a header and a newline respectively
	displaysLine = displaysLine[1 : len(displaysLine)-1]

	for _, displayLine := range displaysLine {
		// For each monitor split by tab and assign
		thisDisplay := Display{Config: DisplayConfig{}}
		columns := strings.Split(displayLine, "\t")

		thisDisplay.ShortMonitorId = columns[15]
		thisDisplay.Adapter = columns[11]
		thisDisplay.Active = columns[3] == "Yes"
		thisDisplay.Disconnected = columns[4] == "Yes"
		thisDisplay.PrimaryMonitor = columns[5] == "Yes"
		thisDisplay.Config.RefreshRate, _ = strconv.Atoi(columns[7])

		// Resolution is in one column. Need to split it apart
		resolution := strings.Split(columns[0], "X")
		thisDisplay.Config.Width, _ = strconv.Atoi(strings.TrimSpace(resolution[0]))
		thisDisplay.Config.Height, _ = strconv.Atoi(strings.TrimSpace(resolution[1]))

		displays = append(displays, thisDisplay)
	}

	return displays
}
