package log

import (
	"log/slog"
	"strings"
)

var Buffer strings.Builder

func Init() *slog.Logger {
	logOpts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&Buffer, &logOpts)
	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
