package infrastructure

import (
	"os"

	"log/slog"
)

func NewLogger() *slog.Logger {
	l := os.Getenv("LOG_LEVEL")

	var ll slog.Level
	switch l {
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
