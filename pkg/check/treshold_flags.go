package check

import (
	"fmt"
	"log/slog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CriticalThreshFlag = "critical"
	WarningThreshFlag  = "warning"
)

var (
	criticalThreshDefault = ""
	warningThreshDefault  = ""
)

const treshDetail = `(e.g. "95% 1s 1000 lookup:1us")
Separate multiple thresholds with space 
Supports %, duration and raw values
To match a single counter, perfix a threshold with a label followed by a colon ':'`

func ThresholdFlags(flags *pflag.FlagSet) {
	flags.String(CriticalThreshFlag, criticalThreshDefault, fmt.Sprintf("Crititcal threshold%s", treshDetail))
	flags.String(WarningThreshFlag, warningThreshDefault, fmt.Sprintf("Warning threshold%s", treshDetail))
}

func SetCriticalThresholdDefault(s string) {
	if len(criticalThreshDefault) > 0 {
		criticalThreshDefault = fmt.Sprintf("%s %s", criticalThreshDefault, s)
	} else {
		criticalThreshDefault = s
	}
	if len(viper.GetString(CriticalThreshFlag)) > 0 {
		slog.Debug("Not setting critical default, parameter given", "param", viper.GetString(CriticalThreshFlag), "default", s)
		return
	}
	slog.Debug("Setting critical default", "default", s)
	viper.SetDefault(CriticalThreshFlag, criticalThreshDefault)
}

func SetWarningThresholdDefault(s string) {
	if len(warningThreshDefault) > 0 {
		warningThreshDefault = fmt.Sprintf("%s %s", warningThreshDefault, s)
	} else {
		warningThreshDefault = s
	}
	if len(viper.GetString(WarningThreshFlag)) > 0 {
		slog.Debug("Not setting warning default, parameter given", "param", viper.GetString(WarningThreshFlag), "default", s)
		return
	}
	slog.Debug("Setting warning default", "default", s)
	viper.SetDefault(WarningThreshFlag, warningThreshDefault)
}
