package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"github.com/maxviazov/excel-flow/internal/config"
)

func New(cfg config.LoggerConfig) (zerolog.Logger, error) {
	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return zerolog.Logger{}, fmt.Errorf("invalid logger config: %w", err)
	}

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("invalid log level %s: %w", cfg.Level, err)
	}
	zerolog.SetGlobalLevel(level)

	var writers []io.Writer

	if cfg.Console {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			NoColor:    !cfg.Color,
		}
		writers = append(writers, consoleWriter)
	}

	if cfg.File != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.File), 0755); err != nil {
			return zerolog.Logger{}, fmt.Errorf("failed to create a log directory: %w", err)
		}
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return zerolog.Logger{}, fmt.Errorf("failed to open log file %s: %w", cfg.File, err)
		}
		writers = append(writers, file)
	}

	var output io.Writer
	if len(writers) == 0 {
		output = os.Stdout
	} else if len(writers) == 1 {
		output = writers[0]
	} else {
		output = io.MultiWriter(writers...)
	}

	logger := zerolog.New(output)

	if cfg.Timestamp {
		logger = logger.With().Timestamp().Logger()
	}

	return logger, nil
}
