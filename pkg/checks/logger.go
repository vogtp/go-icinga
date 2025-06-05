package checks

import (
	"log/slog"
	"strings"
)

var LogBuffer strings.Builder

func initLog() *slog.Logger {
	logOpts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&LogBuffer, &logOpts)
	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
