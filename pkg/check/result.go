package check

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/log"
	"github.com/vogtp/go-icinga/pkg/ssh"
)

type keyValue struct {
	name  string
	value any
}

type Result struct {
	name   string
	prefix string

	header  string
	counter []keyValue
	stati   []keyValue

	counterFormater func(name string, value any) string
	displayFormater func(counter map[string]any) string

	err  error
	code icinga.ResultCode
}

func NewResult(name string, options ...CheckResultOption) *Result {
	r := &Result{
		name:            name,
		stati:           make([]keyValue, 0),
		counter:         make([]keyValue, 0),
		counterFormater: func(name string, value any) string { return fmt.Sprintf("%v", value) },
	}
	for _, o := range options {
		o(r)
	}
	log.Init()
	return r
}

func (r *Result) PrintExit() {
	if len(r.prefix) > 0 && !strings.HasSuffix(r.prefix, ".") {
		r.prefix = fmt.Sprintf("%s.", r.prefix)
	}
	var ret strings.Builder
	var disp strings.Builder
	var pref strings.Builder
	ret.WriteString(r.header)
	if r.err != nil {
		if ret.Len() > 0 {
			ret.WriteString(" - ")
		}
		fmt.Fprintf(&ret, "Error: %v", r.err.Error())
	}
	tw := table.NewWriter()
	style := table.StyleLight
	style.HTML.EscapeText = true
	tw.SetStyle(style)
	for _, c := range r.counter {
		//	pref = fmt.Sprintf("%s%s_ms=%v ", pref, n, t.Milliseconds())
		fmt.Fprintf(&pref, "%s%s=%v ", r.prefix, c.name, r.counterFormater(c.name, c.value))
		//fmt.Fprintf(&disp, "%s\t%v\n", c.name, r.counterFormater(c.name, c.value))
		tw.AppendRow(table.Row{c.name, r.counterFormater(c.name, c.value)})
	}
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignRight},
	})
	disp.WriteString(tw.Render())
	if r.displayFormater != nil {
		ctr := make(map[string]any, len(r.counter))
		for _, c := range r.counter {
			ctr[c.name] = c.value
		}
		disp.Reset()
		disp.WriteString(r.displayFormater(ctr))
	}
	disp.WriteString("\n")
	tw = table.NewWriter()
	tw.SetStyle(style)
	for _, s := range r.stati {
		//fmt.Fprintf(&disp, "%s: %s\n", s.name, s.value)
		tw.AppendRow(table.Row{s.name, s.value})
	}
	disp.WriteString(tw.Render())

	if log.Buffer.Len() > 0 {
		fmt.Fprintf(&disp, "\nLog:\n%s\n", log.Buffer.String())
	}

	isRemote := ssh.IsRemoteRun()

	o := fmt.Sprintf("%s\n\n%s|%s\n", ret.String(), disp.String(), pref.String())

	if isRemote {
		sr := ssh.Result{Out: o, Code: r.code}
		sr.Print()
	} else {
		fmt.Print(o)
	}

	if r.code > icinga.OK {
		os.Exit(int(r.code))
	}
}

func (r *Result) SetHeader(format string, a ...any) {
	r.header = fmt.Sprintf(format, a...)
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
