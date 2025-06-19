package log

import (
	"log/slog"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	verbose = "verbose"
)

var (
	logger *slog.Logger

	Buffer strings.Builder
)

func Init() *slog.Logger {
	if logger != nil {
		return logger
	}
	logOpts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if viper.GetBool(verbose) {
		logOpts.Level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(&Buffer, &logOpts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
	if viper.GetBool(verbose) {
		slog.Debug("Logging verbose")
	}
	return logger
}

func Flags(flags *pflag.FlagSet) {
	flags.Bool(verbose, false, "Log Debug information")
}
