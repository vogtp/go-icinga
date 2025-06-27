package check

import (
	"fmt"

	"github.com/vogtp/go-icinga/pkg/unit"
)

type CheckResultOption func(*Result)

func CheckPrefix(p string) CheckResultOption {
	return func(r *Result) {
		r.prefix = p
	}
}

func DisplayFormater(f func(counter map[string]Data) string) CheckResultOption {
	return func(r *Result) {
		r.displayFormater = f
	}
}
func CounterFormater(f func(name string, d Data) string) CheckResultOption {
	return func(r *Result) {
		r.counterFormater = f
	}
}

func PercentCounterFormater() CheckResultOption {
	return CounterFormater(func(name string, d Data) string {
		if f, ok := d.Value.(float64); ok {
			return fmt.Sprintf("%.3f%%", f)
		}
		return fmt.Sprintf("%v", d.Value)
	},
	)
}

func PercentOrBytesCounterFormater() CheckResultOption {
	return CounterFormater(func(name string, d Data) string {
		if f, ok := d.Value.(float64); ok {
			return fmt.Sprintf("%.3f%%", f)
		}
		if i, ok := d.Value.(uint64); ok {
			return unit.FormatGB(i)
		}
		return fmt.Sprintf("%v", d.Value)
	},
	)
}
