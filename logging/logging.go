package logging

import (
	"log/slog"
	"os"
)

func Initialise(verbose bool) {
	// Create a new text logger with the required log level set
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}
	newDefaultLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(newDefaultLogger)
}
