package powershell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/log"
)

type Session struct {
	cmd *exec.Cmd

	host string
	user string
	pass string

	wg sync.WaitGroup

	stdin       io.WriteCloser
	stdout      bytes.Buffer
	stderr      bytes.Buffer
	debugBuffer strings.Builder
}

func New(ctx context.Context, host string, user string, pass string) (*Session, error) {

	s := &Session{
		host: host,
		user: user,
		pass: pass,
		cmd:  exec.CommandContext(ctx, "pwsh"),
	}
	if viper.GetBool(log.Debug) {
		fmt.Fprintf(&s.debugBuffer, "Started: %s\n", time.Now().Format(time.RFC3339))
	}
	if err := s.init(ctx); err != nil {
		return nil, fmt.Errorf("cannot initalise remote powershell: %w", err)
	}
	s.openRemote()
	return s, nil
}

func (c *Session) Run(ctx context.Context, cmd string) ([]byte, []byte, error) {
	if viper.GetBool(log.Debug) {
		fmt.Fprintln(&c.debugBuffer, cmd)
	}
	c.resetOutput()
	//cmd = fmt.Sprintf("%s%s", remotePath(), cmd)
	c.run(`$out=Invoke-Command -Session $%s -Command { %s }; echo $out`, sessionName, cmd)
	out := c.stdout.String()
	// idx := strings.Index(out, sep)
	// out = out[idx+len(sep):]
	slog.Info("Remote powershell command finished", "cmd", cmd, "stdout", out, "stderr", c.stderr.String())
	return []byte(out), c.stderr.Bytes(), nil
}

func (*Session) CanCopy() bool { return false }

func (c *Session) Copy(ctx context.Context, local, remote string) error {
	return fmt.Errorf("Powershell cannot copy to remote: do it manually") //c.session.Copy(ctx, fmt.Sprintf("%s", local), fmt.Sprintf("%s%s", remotePath(), remote))
}

func (c *Session) Close() {
	if c.cmd == nil {
		return
	}
	c.stdin.Close()
	c.cmd.Cancel()
	slog.Info("Closed powershell session")
	if viper.GetBool(log.Debug) {
		fmt.Fprintf(&c.debugBuffer, "Stopped: %s\n", time.Now().Format(time.RFC3339))
		f, err := os.OpenFile("ps_run.out", os.O_APPEND, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Sprintln(f, c.debugBuffer.String())
		f.Close()
	}
}
