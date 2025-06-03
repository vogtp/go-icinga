package checks

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

type Result struct {
	Name   string
	Prefix string

	Total  time.Duration
	Timing map[string]time.Duration
	Stati  map[string]any

	Err    error
	Result icinga.Result
}

func (r Result) PrintExit() {
	if !strings.HasSuffix(r.Prefix, ".") {
		r.Prefix = fmt.Sprintf("%s.", r.Prefix)
	}
	ret := fmt.Sprintf("%s %s", strings.ToUpper(r.Name), r.Result.String())
	if r.Total > 0 {
		ret = fmt.Sprintf("%s - duration %vµs", ret, r.Total.Microseconds())
	}
	if r.Err != nil {
		ret = fmt.Sprintf("%s - Error: %v", ret, r.Err.Error())
	}
	pref := ""
	disp := ""
	for n, t := range r.Timing {
		//	pref = fmt.Sprintf("%s%s_ms=%v ", pref, n, t.Milliseconds())
		pref = fmt.Sprintf("%s%s%s=%vµs ", pref, r.Prefix, n, t.Microseconds())
		disp = fmt.Sprintf("%s%s\t%vµs\n", disp, n, t.Microseconds())
	}
	for n, s := range r.Stati {
		disp = fmt.Sprintf("%s%s: %s\n", disp, n, s)
	}
	// disp = fmt.Sprintf("%s%s\t%vms", disp, "total", total)
	fmt.Printf("%s\n\n%s | %s", ret, disp, pref)
	if r.Result > icinga.OK {
		os.Exit(int(r.Result))
	}
}
