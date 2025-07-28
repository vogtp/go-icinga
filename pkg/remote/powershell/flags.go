package powershell

import (
	"github.com/spf13/pflag"
)

var IgnoredFlags = []string{RemotingFlag, jeaFlag}

const (
	RemotingFlag = "remote.powershell"
	jeaFlag      = "remote.jea"
)

func Flags(flags *pflag.FlagSet) {
	flags.Bool(RemotingFlag, false, "Use powershell remoting instead of ssh")
	flags.String(jeaFlag, "", "Name of the powershell JEA configuration")
}
