package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/config"
)

// New creates a new zerolog logger with the given configuration
func New(cfg config.LogConfig) zerolog.Logger {
	// Set global log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output
	var output io.Writer = os.Stdout
	if cfg.Output != "stdout" && cfg.Output != "" {
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			output = file
		}
	}

	// Configure error output
	// We're not using separate error output for now, but keeping the config option for future use
	// var errorOutput io.Writer = os.Stderr
	if cfg.ErrorOutput != "stderr" && cfg.ErrorOutput != "" && cfg.ErrorOutput != cfg.Output {
		// Could set up separate error output here if needed
		// For now, we'll just log a message to stdout
		os.Stdout.WriteString("Separate error output configured but not implemented\n")
	}

	// Configure format
	var log zerolog.Logger
	if cfg.Format == "console" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
		log = zerolog.New(output).With().Timestamp().Logger()
	} else {
		// JSON format (default)
		log = zerolog.New(output).With().Timestamp().Logger()
	}

	return log
}
