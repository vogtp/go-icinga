package check

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

// https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/#thresholds

type threshold struct {
	resultCode icinga.ResultCode
	duration   time.Duration
	val        float64
	isPercent  bool
}

func newThreshhold(resultCode icinga.ResultCode, conf string) threshold {
	t := threshold{resultCode: resultCode}
	if strings.HasSuffix(conf, "%") {
		tresh, err := strconv.ParseFloat(conf[:len(conf)-1], 64)
		if err != nil {
			slog.Warn("Cannot parse threshhold", "thresh", tresh, "err", err)
			return t
		}
		t.isPercent = true
		t.val = tresh
		return t
	}
	if f, ok := parseFloat(conf); ok {
		t.val = f
		return t
	}
	d, err := time.ParseDuration(conf)
	if err == nil {
		t.duration = d
		return t
	} else {
		slog.Debug("Cannot parse threshold as duration", "thresh", conf, "err", err)
	}
	return t
}

func (t *threshold) process(kv *keyValue, formatedValue string) icinga.ResultCode {
	kv.resultCode = icinga.OK
	if strings.HasSuffix(formatedValue, "%") {
		if !t.isPercent {
			return icinga.OK
		}
		f, ok := parseFloat(kv.value)
		if !ok {
			slog.Debug("Cannot parse percent float threshhold value", "value", kv.value, "formatedValue", formatedValue)
			return icinga.OK
		}
		if t.val <= f {
			return t.resultCode
		}
		return icinga.OK
	}
	d, err := time.ParseDuration(formatedValue)
	if err == nil {
		if t.duration == 0 {
			return icinga.OK
		}
		if t.duration <= d {
			return t.resultCode
		}
		return icinga.OK
	} else if t.duration != 0 {
		slog.Debug("Cannot parse duration float threshhold value", "value", kv.value, "formatedValue", formatedValue, "err", err)
		return icinga.OK
	}
	f, ok := parseFloat(kv.value)
	if !ok {
		slog.Debug("Cannot parse raw float threshhold value", "value", kv.value, "formatedValue", formatedValue)
		return icinga.OK
	}
	if t.val <= f {
		return t.resultCode
	}
	return icinga.OK
}

func parseFloat(v any) (float64, bool) {
	if f, ok := v.(float64); ok {
		return f, ok
	}
	if f, ok := v.(float32); ok {
		return float64(f), ok
	}
	if i, ok := v.(int); ok {
		return float64(i), ok
	}
	if i, ok := v.(int16); ok {
		return float64(i), ok
	}
	if i, ok := v.(int32); ok {
		return float64(i), ok
	}
	if i, ok := v.(int64); ok {
		return float64(i), ok
	}
	if s, ok := v.(string); ok {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil && len(s) > 0 {
			slog.Debug("Cannot parse threshold float as string", "s", s, "err", err)
			return 0, false
		}
		return f, true
	}
	return 0, false
}
