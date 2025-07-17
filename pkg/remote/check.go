package remote

//go:generate go run gen.go arg1 arg2

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/hash"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/log"
	issh "github.com/vogtp/go-icinga/pkg/ssh"
)

type Session interface {
	Run(cmd string) ([]byte, []byte, error)
	Copy(ctx context.Context, local, remote string) error
}

type client struct {
	user    string
	host    string
	session Session
}

func Check(cmd *cobra.Command, args []string) error {
	log.Init()
	if !ShouldRemoteRun() {
		return nil
	}
	c := client{
		user: viper.GetString(UserFlag),
		host: viper.GetString(HostFlag),
	}
	sess, close, err := issh.New(cmd.Context(), c.user, c.host)
	if err != nil {
		return fmt.Errorf("cannot open ssh session: %w", err)
	}
	defer close()
	c.session = sess
	cmds := strings.Split(cmd.CommandPath(), " ")
	cmds = append(cmds, args...)
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if slices.Contains(ignoredFlags, f.Name) {
			return
		}
		val := f.Value.String()
		if strings.HasSuffix(f.Value.Type(), "Slice") {
			val = strings.ReplaceAll(val, "[", "")
			val = strings.ReplaceAll(val, "]", "")
			val = strings.ReplaceAll(val, ", ", ",")
		}
		//slog.Info("Flag", "name", f.Name, "type", f.Value.Type())
		if f.Value.Type() == "bool" {
			cmds = append(cmds, fmt.Sprintf("--%s", f.Name))
			return
		}
		if len(val) > 0 {
			cmds = append(cmds, fmt.Sprintf("--%s", f.Name), val)
		}
	})

	r, err := c.runRemote(cmd.Context(), cmds)
	if err != nil {
		slog.Info("Remote run error", "err", err, "result", r)
		if log.Buffer.Len() > 0 {
			fmt.Printf("Log:\n%s\n", log.Buffer.String())
		}
	}
	if r.HashMismatch {
		if err := c.copyRemote(cmd.Context(), cmds); err != nil {
			return fmt.Errorf("cannot copy to remote: %w", err)
		}
		r, err = c.runRemote(cmd.Context(), cmds)
		if err != nil {
			return fmt.Errorf("run remote after remote copy: %w", err)
		}
	}
	out := r.Out
	if log.Buffer.Len() > 0 {
		out = strings.ReplaceAll(out, "|", fmt.Sprintf("\nLocal Log:\n%s|", html.EscapeString(log.Buffer.String())))
	}
	if len(out) < 1 {
		fmt.Printf("Log:\n%s\n", log.Buffer.String())
		os.Exit(int(icinga.UNKNOWN))
	}

	fmt.Println(out)
	os.Exit(int(r.Code))
	return nil
}

func (c *client) runRemote(_ context.Context, cmd []string) (*Result, error) {
	if len(cmd) < 1 {
		return nil, fmt.Errorf("no command given: %v", cmd)
	}

	h, err := hash.Calc()
	if err != nil {
		return nil, fmt.Errorf("cannot calculate my hash: %w", err)
	}

	cmdLine := fmt.Sprintf("./%s --%s --%s %q", strings.Join(cmd, " "), isRemoteRun, hashCheckFlag, h)
	slog.Debug("Executing remote command", "cmd", cmdLine, "host", c.host, "user", c.user)
	r, err := c.exec(cmdLine)
	if err != nil {
		return r, fmt.Errorf("%q returned: %w", cmdLine, err)
	}
	return r, nil
}

func (c *client) copyRemote(ctx context.Context, cmd []string) error {

	remote := cmd[0]
	local := os.Args[0]
	slog.Info("remote version is outdated: copy local to remote ", "local", local, "remote", remote)
	if err := c.session.Copy(ctx, local, remote); err != nil {
		return err
	}
	return nil
}

func (c *client) exec(cmd string) (*Result, error) {
	stdo, stde, err := c.session.Run(cmd)
	if len(stde) > 0 {
		fmt.Fprintln(os.Stderr, string(stde))
	}
	if err != nil {
		return nil, err
	}
	r := &Result{}
	if err := json.Unmarshal(stdo, &r); err != nil {
		if log.Buffer.Len() > 0 {
			fmt.Printf("Log:\n%s\n", log.Buffer.String())
		}
		r.HashMismatch = true
		return r, fmt.Errorf("cannot parse remote reponse as json: %w", err)
	}
	return r, err
}
