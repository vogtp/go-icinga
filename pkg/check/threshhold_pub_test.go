package check_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
)

const (
	critTresh = "90% 5s 220"
	warnTresh = "80% 3s 150"
)

func formater(name string, value any) string {
	if strings.HasPrefix(name, "cpu") {
		return fmt.Sprintf("%v%%", value)
	}
	return fmt.Sprintf("%v", value)
}

func TestRaw(t *testing.T) {
	viper.SetDefault(check.CriticalThreshFlag, critTresh)
	viper.Set(check.WarningThreshFlag, warnTresh)
	res := check.NewResult("test", check.CounterFormater(formater))
	tm := check.NewThreshholdsManager(res)
	if rc := tm.Process(); rc != icinga.UNKNOWN {
		t.Errorf("ResultCode is %s not %s", rc, icinga.UNKNOWN)
	}
	res.SetCounter("raw0", 10)
	if rc := tm.Process(); rc != icinga.OK {
		t.Errorf("ResultCode is %s not %s", rc, icinga.OK)
	}
	res.SetCounter("raw1", 150)
	if rc := tm.Process(); rc != icinga.WARNING {
		t.Errorf("ResultCode is %s not %s", rc, icinga.WARNING)
	}
	res.SetCounter("raw2", 220)
	if rc := tm.Process(); rc != icinga.CRITICAL {
		t.Errorf("ResultCode is %s not %s", rc, icinga.CRITICAL)
	}

}
func TestDuration(t *testing.T) {
	viper.Set(check.CriticalThreshFlag, critTresh)
	viper.Set(check.WarningThreshFlag, warnTresh)
	res := check.NewResult("test", check.CounterFormater(formater))
	tm := check.NewThreshholdsManager(res)
	if rc := tm.Process(); rc != icinga.UNKNOWN {
		t.Errorf("ResultCode is %s not %s", rc, icinga.UNKNOWN)
	}
	res.SetCounter("dur0", time.Millisecond*10)
	res.SetCounter("dur1", time.Millisecond*1000)
	if rc := tm.Process(); rc != icinga.OK {
		t.Errorf("ResultCode is %s not %s", rc, icinga.OK)
	}
	res.SetCounter("dur1", time.Millisecond*2999)
	if rc := tm.Process(); rc != icinga.OK {
		t.Errorf("ResultCode is %s not %s", rc, icinga.OK)
	}
	res.SetCounter("dur1", time.Millisecond*3000)
	if rc := tm.Process(); rc != icinga.WARNING {
		t.Errorf("ResultCode is %s not %s", rc, icinga.WARNING)
	}
	res.SetCounter("dur1", time.Second*5)
	if rc := tm.Process(); rc != icinga.CRITICAL {
		t.Errorf("ResultCode is %s not %s", rc, icinga.CRITICAL)
	}
}

func TestPercent(t *testing.T) {
	viper.Set(check.CriticalThreshFlag, critTresh)
	viper.Set(check.WarningThreshFlag, warnTresh)
	res := check.NewResult("test", check.CounterFormater(formater))
	tm := check.NewThreshholdsManager(res)
	if rc := tm.Process(); rc != icinga.UNKNOWN {
		t.Errorf("ResultCode is %s not %s", rc, icinga.UNKNOWN)
	}
	res.SetCounter("cpu0", 10)
	res.SetCounter("cpu1", 50)
	if rc := tm.Process(); rc != icinga.OK {
		t.Errorf("ResultCode is %s not %s", rc, icinga.OK)
	}
	res.SetCounter("cpu4", 80)
	if rc := tm.Process(); rc != icinga.WARNING {
		t.Errorf("ResultCode is %s not %s", rc, icinga.WARNING)
	}
	res.SetCounter("cpu3", 90)
	res.SetCounter("cpu2", 99)
	if rc := tm.Process(); rc != icinga.CRITICAL {
		t.Errorf("ResultCode is %s not %s", rc, icinga.CRITICAL)
	}
	res.SetCounter("cpu5", 0)
	if rc := tm.Process(); rc != icinga.CRITICAL {
		t.Errorf("ResultCode is %s not %s", rc, icinga.CRITICAL)
	}
}

func TestPercentSingle(t *testing.T) {
	viper.Set(check.CriticalThreshFlag, "90%")
	viper.Set(check.WarningThreshFlag, "80%")
	res := check.NewResult("test", check.CounterFormater(formater))
	tm := check.NewThreshholdsManager(res)
	if rc := tm.Process(); rc != icinga.UNKNOWN {
		t.Errorf("ResultCode is %s not %s", rc, icinga.UNKNOWN)
	}
	res.SetCounter("cpu0", 10)
	res.SetCounter("cpu1", 50)
	if rc := tm.Process(); rc != icinga.OK {
		t.Errorf("ResultCode is %s not %s", rc, icinga.OK)
	}
	res.SetCounter("cpu4", 80)
	if rc := tm.Process(); rc != icinga.WARNING {
		t.Errorf("ResultCode is %s not %s", rc, icinga.WARNING)
	}
	res.SetCounter("cpu3", 90)
	res.SetCounter("cpu2", 99)
	if rc := tm.Process(); rc != icinga.CRITICAL {
		t.Errorf("ResultCode is %s not %s", rc, icinga.CRITICAL)
	}
	res.SetCounter("cpu5", 0)
	if rc := tm.Process(); rc != icinga.CRITICAL {
		t.Errorf("ResultCode is %s not %s", rc, icinga.CRITICAL)
	}
}
