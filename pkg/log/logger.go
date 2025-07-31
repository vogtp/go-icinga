package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	Verbose = "verbose"
	Debug   = "debug"
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
		Level:     slog.LevelWarn,
		AddSource: viper.GetBool(Debug),
	}

	var w io.Writer
	w = &Buffer
	if viper.GetBool(Verbose) {
		logOpts.Level = slog.LevelInfo
	}
	if viper.GetBool(Debug) {
		fmt.Fprintln(&Buffer, "******************** THIS BREAKS REMOTE RUNS... ********************")
		logOpts.Level = slog.LevelDebug
		w = io.MultiWriter(&Buffer, os.Stdout)
	}
	handler := slog.NewTextHandler(w, &logOpts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
	if viper.GetBool(Verbose) {
		slog.Debug("Logging verbose")
	}
	return logger
}

func Flags(flags *pflag.FlagSet) {
	flags.Bool(Verbose, false, "Log verbose information")
	flags.Bool(Debug, false, "Log Debug information")
}
