//go:generate stringer -type=ResultCode
package icinga

type ResultCode int

const (
	OK ResultCode = iota
	WARNING
	CRITICAL
	UNKNOWN
)
