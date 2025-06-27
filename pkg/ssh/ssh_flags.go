package ssh

import (
	"log/slog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/director"
)

var ignoredFlags = []string{"help", remoteHost, remoteUser, remoteUserKey, remoteUserKeyPass}

const (
	remoteHost        = "remote.host"
	remoteUser        = "remote.user"
	remoteUserKey     = "remote.sshkey"
	remoteUserKeyPass = "remote.sshkeypass"
	hashCheck         = "hash.check"
	remoteHostDefault = "$host.name$"
	isRemoteRun       = "remote.is_remote"
)

func Flags(flags *pflag.FlagSet, defaultRemoteOn bool) {
	h := ""
	if defaultRemoteOn {
		h = remoteHostDefault
	}
	flags.String(remoteHost, h, "Remote host to run the command on")
	flags.String(remoteUser, "root", "Remote user name")
	flags.String(remoteUserKey, "/var/lib/nagios/.ssh/icinga_ssh", "ssh private key file location")
	flags.String(remoteUserKeyPass, "", "ssh private key password")
	flags.String(hashCheck, "", "check the hash")
	flags.Bool(isRemoteRun, false, "Internal to indicate a remote run")
	if err := flags.MarkHidden(isRemoteRun); err != nil {
		slog.Warn("Cannot hide flag", "flag", isRemoteRun)
	}
	director.IgnoreFlag(isRemoteRun)
	director.IgnoreFlag(hashCheck)
}

// ShouldRemoteRun idicates if the command should be run remotely
func ShouldRemoteRun() bool {
	rh := viper.GetString(remoteHost)
	shouldRunRemote := len(rh) > 0 && rh != remoteHostDefault
	slog.Debug("Should the command run remote", "shouldRunRemote", shouldRunRemote, "remoteHost", rh, "remoteHostDefault", remoteHostDefault)
	return shouldRunRemote
}

func IsRemoteRun() bool {
	return viper.GetBool(isRemoteRun)
}
