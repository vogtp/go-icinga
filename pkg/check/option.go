package check

import "fmt"

type CheckResultOption func(*Result)

func CheckPrefix(p string) CheckResultOption {
	return func(r *Result) {
		r.prefix = p
	}
}

func DisplayFormater(f func(counter map[string]Value) string) CheckResultOption {
	return func(r *Result) {
		r.displayFormater = f
	}
}
func CounterFormater(f func(name string, value Value) string) CheckResultOption {
	return func(r *Result) {
		r.counterFormater = f
	}
}

func PercentCounterFormater() CheckResultOption {
	return CounterFormater(func(name string, value Value) string {
		f, ok := value.Value.(float64)
		if !ok {
			return fmt.Sprintf("%v", value)
		}
		return fmt.Sprintf("%.3f%%", f)
	},
	)
}
