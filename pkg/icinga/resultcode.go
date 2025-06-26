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
	return fmt.Sprintf("[%s]", i.String())
}
