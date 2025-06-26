package check

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

func NewThreshholdsManager(r *Result) *ThreshholdsManager {
	tm := &ThreshholdsManager{
		result:     r,
		thresholds: make([]threshold, 0),
	}
	tm.thresholds = append(tm.thresholds, parseThreshhold(icinga.CRITICAL, CriticalThreshFlag)...)
	tm.thresholds = append(tm.thresholds, parseThreshhold(icinga.WARNING, WarningThreshFlag)...)
	return tm
}

func parseThreshhold(resultCode icinga.ResultCode, flag string) []threshold {
	thres := strings.Split(viper.GetString(flag), " ")
	res := make([]threshold, len(thres))
	for i, t := range thres {
		if len(t) < 1 {
			continue
		}
		res[i] = newThreshhold(resultCode, t)
	}
	slog.Debug("Threshholds loaded", flag, res, "threshholds", thres)
	return res
}

type ThreshholdsManager struct {
	result     *Result
	thresholds []threshold
}

func (tm ThreshholdsManager) Process() icinga.ResultCode {
	resultCode := icinga.UNKNOWN
	for i, c := range tm.result.counter {
		crc := c.resultCode
		for _, t := range tm.thresholds {
			rc := t.process(&c, tm.result.counterFormater(c.name, c.value))
			if resultCode == icinga.UNKNOWN {
				resultCode = rc
			}
			resultCode = max(resultCode, rc)
			crc = max(crc, rc)
		}
		
		tm.result.counter[i].resultCode = crc
	}
	tm.result.code = max(resultCode, tm.result.code)
	return resultCode
}
