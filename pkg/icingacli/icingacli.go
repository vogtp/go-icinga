package icingacli

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func Run(in io.Reader, out io.Writer, args ...string) error {
	cmd := exec.Command("/usr/bin/icingacli", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%q: %w", strings.Join(cmd.Args, " "), err)
	}
	return nil
}
