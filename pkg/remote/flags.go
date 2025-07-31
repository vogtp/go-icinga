package remote

import (
	"log/slog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
	"github.com/vogtp/go-icinga/pkg/log"
	"github.com/vogtp/go-icinga/pkg/remote/powershell"
	"github.com/vogtp/go-icinga/pkg/remote/ssh"
)

var ignoredFlags = []string{"help", log.Debug, HostFlag, UserFlag, passwordFlag, WinRemoteFlag}

const (
	HostFlag      = "remote.host"
	UserFlag      = "remote.user"
	passwordFlag  = "remote.password"
	WinRemoteFlag = "remote.windows"
	RemotePath    = "remote.path"
	hashCheckFlag = "hash.check"
	HostDefault   = "$host.name$"
	isRemoteRun   = "remote.is_remote"
)

func Flags(flags *pflag.FlagSet, defaultRemoteOn bool) {
	h := ""
	if defaultRemoteOn {
		h = HostDefault
	}
	flags.String(HostFlag, h, "Remote host to run the command on")
	flags.String(UserFlag, "root", "Remote user name")
	flags.String(passwordFlag, "", "Remote user password")
	flags.Bool(WinRemoteFlag, false, "Is the remote system a windows system?")
	flags.String(RemotePath, ".", "Remote path for syscheck")
	flags.String(hashCheckFlag, "", "check the hash")
	flags.Bool(isRemoteRun, false, "Internal to indicate a remote run")
	if err := flags.MarkHidden(isRemoteRun); err != nil {
		slog.Warn("Cannot hide flag", "flag", isRemoteRun)
	}
	director.IgnoreFlag(isRemoteRun)
	director.IgnoreFlag(hashCheckFlag)
	ssh.Flags(flags)
	powershell.Flags(flags)
	ignoredFlags = append(ignoredFlags, ssh.IgnoredFlags...)
	ignoredFlags = append(ignoredFlags, powershell.IgnoredFlags...)
}

// ShouldRemoteRun idicates if the command should be run remotely
func ShouldRemoteRun() bool {
	if viper.GetBool(isRemoteRun) {
		return false
	}
	rh := viper.GetString(HostFlag)
	shouldRunRemote := len(rh) > 2 && rh != HostDefault
	slog.Info("Should the command run remote", "shouldRunRemote", shouldRunRemote, "remoteHost", rh, "remoteHostDefault", HostDefault)
	return shouldRunRemote
}

func IsRemoteRun() bool {
	return viper.GetBool(isRemoteRun)
}
