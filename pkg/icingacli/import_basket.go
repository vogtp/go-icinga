package icingacli

import (
	"io"
	"os"
)

func ImportDirectorBasket(r io.Reader) error {
	return Run(r, os.Stdout, "director", "basket", "restore")
}
