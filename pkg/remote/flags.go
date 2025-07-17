package remote

import (
	"log/slog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/go-icinga/pkg/ssh"
)

var ignoredFlags = []string{"help", HostFlag, UserFlag}

const (
	HostFlag        = "remote.host"
	UserFlag        = "remote.user"
	hashCheckFlag         = "hash.check"
	HostDefault     = "$host.name$"
	isRemoteRun           = "remote.is_remote"
)

func Flags(flags *pflag.FlagSet, defaultRemoteOn bool) {
	h := ""
	if defaultRemoteOn {
		h = HostDefault
	}
	flags.String(HostFlag, h, "Remote host to run the command on")
	flags.String(UserFlag, "root", "Remote user name")
	flags.String(hashCheckFlag, "", "check the hash")
	flags.Bool(isRemoteRun, false, "Internal to indicate a remote run")
	if err := flags.MarkHidden(isRemoteRun); err != nil {
		slog.Warn("Cannot hide flag", "flag", isRemoteRun)
	}
	director.IgnoreFlag(isRemoteRun)
	director.IgnoreFlag(hashCheckFlag)
	ssh.Flags(flags)
	ignoredFlags = append(ignoredFlags, ssh.IgnoredFlags...)
}

// ShouldRemoteRun idicates if the command should be run remotely
func ShouldRemoteRun() bool {
	rh := viper.GetString(HostFlag)
	shouldRunRemote := len(rh) > 0 && rh != HostDefault
	slog.Debug("Should the command run remote", "shouldRunRemote", shouldRunRemote, "remoteHost", rh, "remoteHostDefault", HostDefault)
	return shouldRunRemote
}

func IsRemoteRun() bool {
	return viper.GetBool(isRemoteRun)
}
