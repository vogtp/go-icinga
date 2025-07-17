package remote

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/hash"
	"github.com/vogtp/go-icinga/pkg/log"
)

func CheckHash() error {
	log.Init()
	h := viper.GetString(hashCheckFlag)
	if len(h) < 1 {
		return nil
	}
	if err := hash.Check(h); err != nil {
		slog.Info("Remote and local hashes do not match", "err", err)
		r := Result{HashMismatch: true}
		r.Print()
		os.Exit(0)
	}
	return nil
}
