package checks

import (
	"log/slog"
	"strings"
)

var logBuffer strings.Builder

func initLog() *slog.Logger {
	logOpts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&logBuffer, &logOpts)
	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
