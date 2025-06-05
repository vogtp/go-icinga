package checks

import (
	"fmt"
	"os"
	"strings"

	"github.com/vogtp/go-icinga/pkg/icinga"
)

type keyValue struct {
	name  string
	value any
}

type Result struct {
	name   string
	prefix string

	Total   any
	counter []keyValue
	stati   []keyValue

	counterFormater func(name string, value any) string
	displayFormater func(counter map[string]any) string

	err  error
	code icinga.ResultCode
}

func NewCheckResult(name string, options ...CheckResultOption) *Result {
	r := &Result{
		name:            name,
		stati:           make([]keyValue, 0),
		counter:         make([]keyValue, 0),
		counterFormater: func(name string, value any) string { return fmt.Sprintf("%v", value) },
	}
	for _, o := range options {
		o(r)
	}
	return r
}

func (r *Result) PrintExit() {
	if len(r.prefix) > 0 && !strings.HasSuffix(r.prefix, ".") {
		r.prefix = fmt.Sprintf("%s.", r.prefix)
	}
	var ret strings.Builder
	var disp strings.Builder
	var pref strings.Builder
	ret.WriteString(r.code.String())
	if r.Total != nil {
		fmt.Fprintf(&ret, " - total %v", r.counterFormater("total", r.Total))
	}
	if r.err != nil {
		fmt.Fprintf(&ret, " - Error: %v", r.err.Error())
	}
	for _, c := range r.counter {
		//	pref = fmt.Sprintf("%s%s_ms=%v ", pref, n, t.Milliseconds())
		fmt.Fprintf(&pref, "%s%s=%v ", r.prefix, c.name, r.counterFormater(c.name, c.value))
		fmt.Fprintf(&disp, "%s\t%v\n", c.name, r.counterFormater(c.name, c.value))
	}
	if r.displayFormater != nil {
		ctr := make(map[string]any, len(r.counter))
		for _, c := range r.counter {
			ctr[c.name] = c.value
		}
		disp.Reset()
		disp.WriteString(r.displayFormater(ctr))
	}
	for _, s := range r.stati {
		fmt.Fprintf(&disp, "%s: %s\n", s.name, s.value)
	}

	if LogBuffer.Len() > 0 {
		fmt.Fprintf(&disp, "\nLog:\n%s\n", LogBuffer.String())
	}

	fmt.Printf("%s\n\n%s|%s", ret.String(), disp.String(), pref.String())
	if r.code > icinga.OK {
		os.Exit(int(r.code))
	}
}

func (r *Result) SetCode(c icinga.ResultCode) {
	r.code = max(r.code, c)
}

func (r *Result) SetCounter(name string, val any) {
	r.counter = append(r.counter, keyValue{name: name, value: val})
}

func (r *Result) SetStatus(name string, val any) {
	r.stati = append(r.stati, keyValue{name: name, value: val})
}

func (r *Result) SetError(err error) {
	r.err = err
	r.SetCode(icinga.WARNING)
}
