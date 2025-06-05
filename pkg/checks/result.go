package checks

import (
	"fmt"
	"os"
	"strings"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

type Result struct {
	name   string
	prefix string

	Total   any
	counter map[string]any
	stati   map[string]any

	counterFormater func(name string, value any) string
	displayFormater func(counter map[string]any) string

	err  error
	code icinga.ResultCode
}

func NewCheckResult(name string, options ...CheckResultOption) *Result {
	r := &Result{
		name:            name,
		stati:           make(map[string]any),
		counter:         make(map[string]any),
		counterFormater: func(name string, value any) string { return fmt.Sprintf("%v", value) },
	}
	for _, o := range options {
		o(r)
	}
	return r
}

func (r Result) PrintExit() {
	if len(r.prefix) > 0 && !strings.HasSuffix(r.prefix, ".") {
		r.prefix = fmt.Sprintf("%s.", r.prefix)
	}
	ret := r.code.String()
	if r.Total != nil {
		ret = fmt.Sprintf("%s - total %v", ret, r.counterFormater("total", r.Total))
	}
	if r.err != nil {
		ret = fmt.Sprintf("%s - Error: %v", ret, r.err.Error())
	}
	pref := ""
	disp := ""
	for n, c := range r.counter {
		//	pref = fmt.Sprintf("%s%s_ms=%v ", pref, n, t.Milliseconds())
		pref = fmt.Sprintf("%s%s%s=%v ", pref, r.prefix, n, r.counterFormater(n, c))
		disp = fmt.Sprintf("%s%s\t%v\n", disp, n, r.counterFormater(n, c))
	}
	if r.displayFormater != nil {
		disp = r.displayFormater(r.counter)
	}
	for n, s := range r.stati {
		disp = fmt.Sprintf("%s%s: %s\n", disp, n, s)
	}

	if LogBuffer.Len() > 0 {
		disp = fmt.Sprintf("%s\nLog:\n%s\n", disp, LogBuffer.String())
	}

	fmt.Printf("%s\n\n%s | %s", ret, disp, pref)
	if r.code > icinga.OK {
		os.Exit(int(r.code))
	}
}

func (r *Result) SetCode(c icinga.ResultCode) {
	r.code = max(r.code, c)
}

func (r *Result) SetCounter(name string, val any) {
	r.counter[name] = val
}

func (r *Result) SetStatus(name string, val any) {
	r.stati[name] = val
}

func (r *Result)SetError(err error){
	r.err = err
	r.SetCode(icinga.WARNING)
}