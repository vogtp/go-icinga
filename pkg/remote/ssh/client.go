package ssh

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

type Session struct {
	sshClient *ssh.Client
}

func New(ctx context.Context, host string, user string, pass string) (*Session, error) {

	authMethods, err := getSshAuth()
	if err != nil {
		return nil, fmt.Errorf("no ssh auth: %w", err)
	}
	if len(pass) > 0 {
		authMethods = append(authMethods, ssh.Password(pass))
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	c := &Session{
		sshClient: client,
	}
	return c, nil
}

func (c *Session) Run(cmd string) ([]byte, []byte, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	var stdo bytes.Buffer
	var stde bytes.Buffer
	session.Stdout = &stdo
	session.Stderr = &stde
	err = session.Run(cmd)
	return stdo.Bytes(), stde.Bytes(), err
}

func (c *Session) Copy(ctx context.Context, local, remote string) error {
	client, err := scp.NewClientBySSH(c.sshClient)
	if err != nil {
		return fmt.Errorf("error creating new SSH session from existing connection: %w", err)
	}
	defer client.Close()

	f, err := os.Open(local)
	if err != nil {
		return fmt.Errorf("error opening local file %q: %w", local, err)
	}
	defer f.Close()

	if err := client.CopyFromFile(ctx, *f, remote, "0755"); err != nil {
		return fmt.Errorf("error while scp file %q: %w", remote, err)
	}
	return nil
}

func (c *Session) Close() {
	if c.sshClient == nil {
		return
	}
	if err := c.sshClient.Close(); err != nil {
		slog.Warn("Error closing ssh client", "err", err)
	}
}
