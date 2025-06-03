//go:generate stringer -type=Result
package icinga

type Result int

const (
	OK Result = iota
	WARNING
	CRITICAL
	UNKNOWN
)