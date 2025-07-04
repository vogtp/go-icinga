package check

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/go-icinga/pkg/log"
	"github.com/vogtp/go-icinga/pkg/ssh"
)

type Data struct {
	Name              string
	Value             any
	ResultCode        icinga.ResultCode
	CriticalThreshold any
	WaningThreshold   any
}

type Result struct {
	name   string
	prefix string

	header  string
	counter []Data
	stati   []Data

	counterFormater func(name string, value Data) string
	displayFormater func(counter map[string]Data) string

	err  error
	code icinga.ResultCode
}

func NewResult(name string, options ...CheckResultOption) *Result {
	r := &Result{
		name:            name,
		stati:           make([]Data, 0),
		counter:         make([]Data, 0),
		counterFormater: func(name string, value Data) string { return fmt.Sprintf("%v", value.Value) },
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
	tm := NewThreshholdsManager(r)
	tm.Process()
	for _, c := range r.counter {
		fmtCnt := r.counterFormater(c.Name, c)
		fmt.Fprintf(&pref, "'%s%s'=%v%v ", r.prefix, c.Name, fmtCnt, getPrefDataThreshDisplay(c))
		tw.AppendRow(table.Row{c.ResultCode.IcingaString(), c.Name, fmtCnt})
	}
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Align: text.AlignRight},
	})
	disp.WriteString(tw.Render())
	if r.displayFormater != nil {
		ctr := make(map[string]Data, len(r.counter))
		for _, c := range r.counter {
			ctr[c.Name] = c
		}
		disp.Reset()
		disp.WriteString(r.displayFormater(ctr))
	}
	disp.WriteString("\n")
	tw = table.NewWriter()
	tw.SetStyle(style)
	for _, s := range r.stati {
		//fmt.Fprintf(&disp, "%s: %s\n", s.name, s.value)
		tw.AppendRow(table.Row{s.Name, s.Value})
	}
	disp.WriteString(tw.Render())

	if log.Buffer.Len() > 0 {
		fmt.Fprintf(&disp, "\nLog:\n%s\n", log.Buffer.String())
	}

	isRemote := ssh.IsRemoteRun()
	slog.Debug("Is this command running by ssh?", "isRemote", isRemote)
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

func getPrefDataThreshDisplay(d Data) string {
	if d.WaningThreshold == nil && d.CriticalThreshold == nil {
		return ""
	}
	w := d.WaningThreshold
	if w == nil {
		w = ""
	}
	c := d.CriticalThreshold
	if c == nil {
		c = ""
	}
	return fmt.Sprintf(";%v;%v", w, c)
}

func (r *Result) SetHeader(format string, a ...any) {
	r.header = fmt.Sprintf(format, a...)
}

func (r *Result) GetCode() icinga.ResultCode {
	return r.code
}

func (r *Result) SetCode(c icinga.ResultCode) {
	r.code = max(r.code, c)
}

func (r *Result) SetCounter(name string, val any) {
	r.counter = append(r.counter, Data{Name: name, Value: val})
}

func (r *Result) SetStatus(name string, val any) {
	r.stati = append(r.stati, Data{Name: name, Value: val})
}

func (r *Result) SetError(err error) {
	r.err = err
	r.SetCode(icinga.WARNING)
}

func (d Data) String() string {
	return fmt.Sprintf("%v", d.Value)
}

func (d *Data) SetThreshold(t *threshold) {
	rt := reflect.TypeOf(d.Value).String()
	var v any
	if strings.HasPrefix(rt, "uint") || strings.HasPrefix(rt, "int") {
		v = int(t.val)
	} else if rt == "time.Duration" {
		v = t.duration.Microseconds()
	} else {
		// fmt.Printf("type %v\n", rt)
		v = t.val
	}
	if t.resultCode == icinga.CRITICAL {
		d.CriticalThreshold = v
	} else {
		d.WaningThreshold = v
	}
}
