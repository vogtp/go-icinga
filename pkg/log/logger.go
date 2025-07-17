package log

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	verbose = "verbose"
	debug   = "debug"
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
		Level: slog.LevelWarn,
	}
	var w io.Writer
	w = &Buffer
	if viper.GetBool(verbose) {
		logOpts.Level = slog.LevelInfo
	}
	if viper.GetBool(debug) {
		logOpts.Level = slog.LevelDebug
		w = io.MultiWriter(&Buffer, os.Stdout)
	}
	handler := slog.NewTextHandler(w, &logOpts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
	if viper.GetBool(verbose) {
		slog.Debug("Logging verbose")
	}
	return logger
}

func Flags(flags *pflag.FlagSet) {
	flags.Bool(verbose, false, "Log verbose information")
	flags.Bool(debug, false, "Log Debug information")
}
