package remote

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

type Result struct {
	Out          string
	HashMismatch bool
	Code         icinga.ResultCode
}

func (r Result) Print() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	if err := enc.Encode(r); err != nil {
		slog.Warn("Cannot print remote ssh response", "err", err)
	}
}
