package logger

import (
	"github.com/jesses-code-adventures/treeai/config"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init(cfg *config.Config) {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	if cfg.Silent {
		level = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	// Add config as context to all log entries
	Logger = slog.New(handler).With(cfg.ToSlogAttrs()...)
}
