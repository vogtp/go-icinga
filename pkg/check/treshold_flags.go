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

func ThresholdFlags(flags *pflag.FlagSet) {
	flags.String(CriticalThreshFlag, criticalThreshDefault, "Crititcal threshold, separate multiple thresholds with space, supports duration, % and raw values ")
	flags.String(WarningThreshFlag, warningThreshDefault, "Warning threshold, separate multiple thresholds with space, supports duration, % and raw values")
}

func SetCriticalThresholdDefault(flags *pflag.FlagSet, s string) {
	if len(criticalThreshDefault) > 0 {
		criticalThreshDefault = fmt.Sprintf("%s %s", criticalThreshDefault, s)
	} else {
		criticalThreshDefault = s
	}
	if len(viper.GetString(CriticalThreshFlag)) > 0 {
		slog.Debug("Not setting critical default, parameter given", "param", viper.GetString(CriticalThreshFlag), "default", s)
		return
	}
	slog.Debug("Setting critical default", "param", viper.GetString(CriticalThreshFlag), "default", s)
	viper.Set(CriticalThreshFlag, criticalThreshDefault)
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
	slog.Debug("Setting warning default", "param", viper.GetString(WarningThreshFlag), "default", s)
	viper.Set(WarningThreshFlag, warningThreshDefault)
}
