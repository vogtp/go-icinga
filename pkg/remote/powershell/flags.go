package powershell

import (
	"github.com/spf13/pflag"
)

var IgnoredFlags = []string{RemotingFlag,  JeaFlag}

const (
	RemotingFlag = "remote.powershell"
	JeaFlag        = "remote.jea"
)

func Flags(flags *pflag.FlagSet) {
	flags.Bool(RemotingFlag, false, "Use powershell remoting instead of ssh")
	flags.String(JeaFlag, "", "Name of the powershell JEA configuration")
}
