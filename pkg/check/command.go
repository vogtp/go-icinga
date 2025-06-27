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
	goicinga "github.com/vogtp/go-icinga"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/log"
	"github.com/vogtp/go-icinga/pkg/ssh"
)

type Command struct {
	*cobra.Command
	Use             string
	Short           string
	NamePrefix      string // director prefix
	DescriptionURL  string // URL for diectory config
	Criticality     icinga.Criticality
	DefaultRemoteOn bool // should we run on the remote host by default

	preRun func(cmd *cobra.Command, args []string) error
}

func (c *Command) Execute() error {
	return c.ExecuteContext(context.Background())
}

const (
	versionFlag = "version"
)

func (c *Command) ExecuteContext(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	c.init()

	flags := c.PersistentFlags()
	log.Flags(flags)
	ssh.Flags(flags, c.DefaultRemoteOn)
	ThresholdFlags(flags)
	director.Flags(flags)
	flags.Bool(versionFlag, false, "Prints the version")
	flags.VisitAll(func(f *pflag.Flag) {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(err)
		}
	})
	if c.PersistentPreRunE != nil {
		c.preRun = c.PersistentPreRunE
	}
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if viper.GetBool(versionFlag) {
			fmt.Println(goicinga.Version())
			os.Exit(0)
		}
		if c.preRun != nil {
			if err := c.preRun(cmd, args); err != nil {
				return err
			}
		}
		if cmd.PreRunE != nil {
			if err := cmd.PreRunE(cmd, args); err != nil {
				return err
			}
		}
		if err := ssh.CheckHash(); err != nil {
			return fmt.Errorf("cannot check hash: %w", err)
		}
		if err := c.generateDirectorConfig(cmd, args); err != nil {
			fmt.Printf("Icinga director import error: %v\n", err)
			os.Exit(1)
			return nil
		}
		//	fmt.Printf("ssh key: %s\n", viper.GetString("remote.sshkey"))
		if err := ssh.RemoteCheck(cmd, args); err != nil {
			u, err2 := user.Current()
			slog.Warn("Remote check error", "username", u.Name, "home", u.HomeDir, "errr", err2)
			fmt.Printf("Remote run returned error: %v\n", err)
			os.Exit(int(icinga.UNKNOWN))
			return nil
		}
		return nil
	}
	if c.RunE == nil {
		c.RunE = func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		}
	}

	return c.Command.ExecuteContext(ctx)
}

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmds ...*cobra.Command) {
	c.init()
	c.Command.AddCommand(cmds...)
}

func (c *Command) generateDirectorConfig(cmd *cobra.Command, _ []string) error {
	if !director.ShouldGenerate() {
		return nil
	}
	slog.Info("thresh", "crit", viper.GetString("critical"))
	slog.Debug("Generate dir config")
	d := director.Generator{
		NamePrefix:     c.NamePrefix,
		Description:    c.Use,
		DescriptionURL: c.DescriptionURL,
		CobraCmd:       cmd,
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
	return nil
}
func (c *Command) init() {
	if c.Command == nil {
		c.Command = &cobra.Command{}
	}
	// wirte those values too if the command was generated
	c.Command.Use = c.Use
	c.Command.Short = c.Short
}
