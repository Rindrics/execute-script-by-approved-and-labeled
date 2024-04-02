package infrastructure

import (
	"os"

	"log/slog"

	"github.com/Rindrics/execute-script-with-merge/domain"
)

func NewLogger() *slog.Logger {
	l := os.Getenv(domain.EnvVarLogLevel)

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
