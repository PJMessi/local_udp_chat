package logger

import (
	"log/slog"
	"os"
)

const LOG_LEVEL = slog.LevelInfo

func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: LOG_LEVEL,
	}))

	return logger
}
