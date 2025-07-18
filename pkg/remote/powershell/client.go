package powershell

import (
	"context"
	"fmt"
	"log/slog"

	ps "github.com/vogtp/go-powershell"
	"github.com/vogtp/go-powershell/backend"
	"github.com/vogtp/go-powershell/middleware"
)

type Session struct {
	shell   ps.Shell
	session middleware.Middleware
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

func (c *Session) Run(cmd string) ([]byte, []byte, error) {
	stdout, stderr, err := c.session.Execute(cmd)
	if err != nil {
		return []byte(stdout), []byte(stderr), err
	}
	return []byte(stdout), []byte(stderr), nil
}

func (c *Session) Copy(ctx context.Context, local, remote string) error {
	return c.session.Copy(fmt.Sprintf("%s.exe", local), fmt.Sprintf("c:\\%s.exe", remote))
}

func (c *Session) Close() {
	if c.session == nil {
		return
	}
	c.session.Exit()
	slog.Info("Closed powershell session")
}
