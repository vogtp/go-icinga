package checks

import (
	"fmt"
	"os"
	"strings"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

type Result struct {
	Name   string
	Prefix string

	Total   any
	Counter map[string]any
	Stati   map[string]any

	CounterFormater func(name string, value any) string

	Err    error
	Result icinga.Result
}

func (r Result) PrintExit() {
	if r.CounterFormater == nil {
		r.CounterFormater = func(name string, value any) string { return fmt.Sprintf("%v", value) }
	}
	if len(r.Prefix) > 0 && !strings.HasSuffix(r.Prefix, ".") {
		r.Prefix = fmt.Sprintf("%s.", r.Prefix)
	}
	ret := fmt.Sprintf("%s %s", strings.ToUpper(r.Name), r.Result.String())
	if r.Total != nil {
		ret = fmt.Sprintf("%s - total %v", ret, r.CounterFormater("total", r.Total))
	}
	if r.Err != nil {
		ret = fmt.Sprintf("%s - Error: %v", ret, r.Err.Error())
	}
	pref := ""
	disp := ""
	for n, c := range r.Counter {
		//	pref = fmt.Sprintf("%s%s_ms=%v ", pref, n, t.Milliseconds())
		pref = fmt.Sprintf("%s%s%s=%v ", pref, r.Prefix, n, r.CounterFormater(n, c))
		disp = fmt.Sprintf("%s%s\t%v\n", disp, n, r.CounterFormater(n, c))
	}
	for n, s := range r.Stati {
		disp = fmt.Sprintf("%s%s: %s\n", disp, n, s)
	}

	if logBuffer.Len() > 0 {
		disp = fmt.Sprintf("%s\nLog:\n%s\n", disp, logBuffer.String())
	}

	fmt.Printf("%s\n\n%s | %s", ret, disp, pref)
	if r.Result > icinga.OK {
		os.Exit(int(r.Result))
	}
}
