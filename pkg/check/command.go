package check

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/ssh"
)

type Command struct {
	*cobra.Command
	Use            string
	Short          string
	NamePrefix     string
	DescriptionURL string
	Criticality    icinga.Criticality

	preRun func(cmd *cobra.Command, args []string) error
}

func (c *Command) Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	c.init()

	flags := c.PersistentFlags()
	ssh.Flags(flags)
	director.Flags(flags)
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	if c.PersistentPreRunE != nil {
		c.preRun = c.PersistentPreRunE
	}
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if c.preRun != nil {
			if err := c.preRun(cmd, args); err != nil {
				return err
			}
		}
		if err := c.generateDirectorConfig(args); err != nil {
			return err
		}
		//	fmt.Printf("ssh key: %s\n", viper.GetString("remote.sshkey"))
		if err := ssh.RemoteCheck(cmd, args); err != nil {
			u, err2 := user.Current()
			slog.Warn("Remote check error", "username", u.Name, "home", u.HomeDir, "errr", err2)
			return err
		}
		return nil
	}
	if c.RunE == nil {
		c.RunE = func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		}
	}

	return c.ExecuteContext(ctx)
}

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmds ...*cobra.Command) {
	c.init()
	c.Command.AddCommand(cmds...)
}

func (c *Command) generateDirectorConfig(args []string) error {
	if director.ShouldGenerate() {
		d := director.Generator{
			NamePrefix:     c.NamePrefix,
			Description:    c.Use,
			DescriptionURL: c.DescriptionURL,
			CobraCmd:       c.Command,
			Output:         os.Stdout,
			Criticality:    c.Criticality,
		}
		if len(c.Short) > 0 {
			d.Description = fmt.Sprintf("%s: %s", d.Description, c.Short)
		}
		if err := d.Generate(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}
func (c *Command) init() {
	if c.Command == nil {
		c.Command = &cobra.Command{
			Use:   c.Use,
			Short: c.Short,
		}
	}
}
