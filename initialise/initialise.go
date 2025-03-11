package initialise

import (
	"streamres/update"
)

// Tool makes sure the dependencies are initialised into the cache directory before running anything.
// This tool relies on several executables (see readme) to work so needs them places in the user cache dir
func Tool() error {
	return update.Tool()
}
