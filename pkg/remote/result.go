package remote

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

type Result struct {
	Out          []byte
	HashMismatch bool
	Code         icinga.ResultCode
}

func (r Result) Print() {
	// b64 := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	// defer b64.Close()
	enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", " ")
	if err := enc.Encode(r); err != nil {
		slog.Warn("Cannot print remote ssh response", "err", err)
	}
}
