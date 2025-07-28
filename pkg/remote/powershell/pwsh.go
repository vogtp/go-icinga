package powershell

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

func (c *Session) init(ctx context.Context) error {
	c.cmd.Env = append(c.cmd.Environ(), "KRB5_CONFIG=/etc/krb5_ADS.conf")
	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return err
	}
	go c.handleOut(ctx, "stdout", stdout, &c.stdout)
	go c.handleOut(ctx, "stderr", stderr, &c.stdout)
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	c.stdin = stdin
	if err := c.cmd.Start(); err != nil {
		return err
	}
	c.run(`function prompt {"%s"}`, prompt)
	c.run("$PSStyle.OutputRendering = [System.Management.Automation.OutputRendering]::PlainText;")

	return nil
}

func (c *Session) openRemote() {
	c.run(" echo %s | kinit %s", c.pass, c.user)
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		c.run("klist -l")
	}
	jeaConfigName := viper.GetString(JeaFlag)
	if len(jeaConfigName) > 0 {
		jeaConfigName = fmt.Sprintf("-ConfigurationName '%s'", jeaConfigName)
		slog.Info("Using JEA", "configuration", jeaConfigName)
	}
	c.run("$%s = New-PSSession %s %s", sessionName, c.host, jeaConfigName)
	c.resetOutput()
}

const (
	sessionName = "GoPowershellRemoteSessionName"
	endOfCmd    = "GoPowershellCommandFinishSign"
	prompt      = "GOPS"
)

func (c *Session) Wait() {
	c.wg.Wait()
}
func (c *Session) resetOutput() {
	c.stdout.Reset()
	c.stderr.Reset()
}

func (c *Session) handleOut(ctx context.Context, name string, r io.Reader, w io.Writer) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		t := s.Text()
		//for _, t := range strings.Split(s.Text(), "")
		if strings.HasPrefix(t, prompt) {
			continue
		}
		if strings.Contains(t, endOfCmd) {
			if strings.Contains(t, "echo") {
				continue
			}
			slog.Debug("End of command found", "text", t)
			c.wg.Done()
			continue
		}
		if _, err := w.Write([]byte(t)); err != nil {
			slog.Warn("Cannot write powershell output", "writer", name, "err", err)
		}
		slog.Debug("Remote powershell", name, c.redactPassword(t))
		if ctx.Err() != nil {
			break
		}
	}
}

func (s *Session) run(cmd string, a ...any) {
	s.wg.Add(1)
	cmd = fmt.Sprintf(cmd, a...)
	s.debugBuffer.WriteString(cmd)
	fmt.Fprintf(s.stdin, "%s\n", cmd)
	fmt.Fprintf(s.stdin, "echo '%s'\n", endOfCmd)

	slog.Debug("Remote powershell", "stdin", s.redactPassword(cmd))
	s.wg.Wait()
}

func (s *Session) redactPassword(str string) string {
	if strings.Contains(str, s.pass) {
		return strings.ReplaceAll(str, s.pass, "<REDACTED PASSWORD>")
	}
	return str
}
