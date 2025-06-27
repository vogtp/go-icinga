//go:generate stringer -type=ResultCode
package icinga

import "fmt"

type ResultCode int

const (
	OK ResultCode = iota
	WARNING
	CRITICAL
	UNKNOWN
)

func (i ResultCode) IcingaString() string {
	switch i {
	case OK:
		return "🟢"
	case WARNING:
		return "🟠"
	case CRITICAL:
		return "🔴"
	case UNKNOWN:
		return "🔵"
	}
	return fmt.Sprintf("[%s]", i.String())
}
