package ssh

//go:generate go run gen.go arg1 arg2

import (
	"bytes"
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
	"golang.org/x/crypto/ssh"
)

func RemoteCheck(cmd *cobra.Command, args []string) error {
	log.Init()
	if !ShouldRemoteRun() {
		return nil
	}
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

	r, err := runRemote(cmd.Context(), cmds)
	if err != nil {
		slog.Info("Remote error", "err", err, "result", r)
		if log.Buffer.Len() > 0 {
			fmt.Printf("Log:\n%s\n", log.Buffer.String())
		}
	}
	if r.HashMismatch {
		if err:=copyRemote(cmd.Context(), cmds); err != nil{
			return fmt.Errorf("cannot copy to remote: %w",err)
		}
		_, err := runRemote(cmd.Context(), cmds)
		return err
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

func CheckHash() error {
	log.Init()
	h := viper.GetString(hashCheck)
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

func runRemote(ctx context.Context, cmd []string) (*Result, error) {
	if len(cmd) < 1 {
		return nil, fmt.Errorf("no command given: %v", cmd)
	}
	user := viper.GetString(remoteUser)
	host := viper.GetString(remoteHost)

	sshAuth, err := getSshAuth()
	if err != nil {
		return nil, fmt.Errorf("no ssh auth: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            sshAuth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	h, err := hash.Calc()
	if err != nil {
		return nil, fmt.Errorf("cannot calculate my hash: %w", err)
	}

	cmdLine := fmt.Sprintf("./%s --%s --%s %q", strings.Join(cmd, " "), isRemoteRun, hashCheck, h)
	slog.Debug("Executing remote command", "cmd", cmdLine, "host", host, "user", user)
	r, err := exec(client, cmdLine)
	if err != nil {
		return r, fmt.Errorf("%q returned: %w", cmdLine, err)
	}
	return r, nil
}

func copyRemote(ctx context.Context, cmd []string)  error {
	user := viper.GetString(remoteUser)
	host := viper.GetString(remoteHost)

	sshAuth, err := getSshAuth()
	if err != nil {
		return fmt.Errorf("no ssh auth: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            sshAuth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return  fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()
	remote := cmd[0]
	local := os.Args[0]
	slog.Info("remote version is outdated: copy local to remote ", "local", local, "remote", remote)
	if err := Copy(ctx, client, local, remote); err != nil {
		return  err
	}
	return nil
}

func exec(client *ssh.Client, cmd string) (*Result, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	var stdo bytes.Buffer
	var stde bytes.Buffer
	session.Stdout = &stdo
	session.Stderr = &stde
	err = session.Run(cmd)
	if stde.Len() > 0 {
		fmt.Println(stde.String())
	}
	r := &Result{}
	if err := json.Unmarshal(stdo.Bytes(), &r); err != nil {
		if log.Buffer.Len() > 0 {
			fmt.Printf("Log:\n%s\n", log.Buffer.String())
		}
		r.HashMismatch = true
		return r, fmt.Errorf("cannot parse remote reponse as json: %w", err)
	}
	return r, err
}
