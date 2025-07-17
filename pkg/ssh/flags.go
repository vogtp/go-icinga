package ssh

import (
	"github.com/spf13/pflag"
)

var IgnoredFlags = []string{UserKeyFlag, UserKeyPassFlag}

const (
	UserKeyFlag     = "remote.sshkey"
	UserKeyPassFlag = "remote.sshkeypass"
)

func Flags(flags *pflag.FlagSet) {
	flags.String(UserKeyFlag, "/var/lib/nagios/.ssh/icinga_ssh", "ssh private key file location")
	flags.String(UserKeyPassFlag, "", "ssh private key password")
}
