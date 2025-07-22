package powershell

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/log"
	ps "github.com/vogtp/go-powershell"
	"github.com/vogtp/go-powershell/backend"
	"github.com/vogtp/go-powershell/middleware"
)

type Session struct {
	shell       ps.Shell
	session     middleware.Middleware
	debugBuffer strings.Builder
}

func remotePath() string {
	return `C:\Program Files\WindowsPowerShell\Modules\UNIBAS.JEA.WSH\0.5.0\Public\`
}

func New(ctx context.Context, host string, user string, pass string) (*Session, error) {

	back := &backend.Local{}
	shell, err := ps.New(back)
	if err != nil {
		return nil, err
	}
	// test(shell, `$a="hello"; $a`)
	// test(shell, `$a`)
	// test(shell, `$n="world"; "$a $n"`)
	// test(shell, `"$a $n"`)
	config := middleware.NewSessionConfig()
	config.ComputerName = host
	config.Credential = &middleware.UserPasswordCredential{
		Username: user,
		Password: pass,
	}
	session, err := middleware.NewSession(shell, config)
	if err != nil {
		return nil, err
	}
	s := &Session{
		shell:   shell,
		session: session,
	}
	if viper.GetBool(log.Debug) {
		fmt.Fprintf(&s.debugBuffer, "Started: %s\n", time.Now().Format(time.RFC3339))
	}
	return s, nil
}

// func test(shell ps.Shell, cmd string) {
// 	fmt.Printf("starting cmd:%s\n", cmd)
// 	o, e, err := shell.Execute(cmd)
// 	fmt.Printf("cmd:%s\n", cmd)
// 	if err != nil {

// 		fmt.Printf("err: %v\n\n", err.Error())
// 		return
// 	}
// 	fmt.Printf("OUT:%s\nErr:%s\n\n", o, e)
// }

func (c *Session) Run(ctx context.Context, cmd string) ([]byte, []byte, error) {
	if viper.GetBool(log.Debug) {
		fmt.Fprintln(&c.debugBuffer, cmd)
	}
	stdout, stderr, err := c.session.Execute(ctx, fmt.Sprintf("%s%s", remotePath(), cmd))
	if err != nil {
		return []byte(stdout), []byte(stderr), err
	}
	return []byte(stdout), []byte(stderr), nil
}

func (c *Session) Copy(ctx context.Context, local, remote string) error {
	return c.session.Copy(ctx, fmt.Sprintf("%s", local), fmt.Sprintf("%s%s", remotePath(), remote))
}

func (c *Session) Close() {
	if c.session == nil {
		return
	}
	c.session.Exit()
	slog.Info("Closed powershell session")
	if viper.GetBool(log.Debug) {
		fmt.Fprintf(&c.debugBuffer, "Stopped: %s\n", time.Now().Format(time.RFC3339))
		f, err := os.OpenFile("ps_run.out", os.O_APPEND, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Sprintln(f, c.debugBuffer.String())
	}
}
