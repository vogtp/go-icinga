package powershell

import (
	"github.com/spf13/pflag"
)

var IgnoredFlags = []string{RemotingFlag, psConfig, krbcfgFlag}

const (
	RemotingFlag = "remote.powershell"
	psConfig     = "remote.powershell.config"
	krbcfgFlag   = "remote.krb5_config"
)

func Flags(flags *pflag.FlagSet) {
	flags.Bool(RemotingFlag, false, "Use powershell remoting instead of ssh")
	flags.String(psConfig, "", "Name of the powershell (JEA) configuration")
	flags.String(krbcfgFlag, "/etc/krb5.conf", "Kerberos config file")
}
