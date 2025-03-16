package ctxutils

import (
	"context"
	"log/slog"
	"os"
)

type loggerKey string

var key loggerKey = "base_logger"

func AttachLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func ExtractLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(key).(*slog.Logger); ok {
		return logger
	}

	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
