package infrastructure

import (
	"os"

	"log/slog"
)

func NewLogger(loglevel string) *slog.Logger {
	var ll slog.Level
	switch loglevel {
	case "debug":
		ll = slog.LevelDebug
	case "info":
		ll = slog.LevelInfo
	case "warn":
		ll = slog.LevelWarn
	case "error":
		ll = slog.LevelError
	default:
		ll = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &ll}))
}
