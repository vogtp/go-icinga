package icinga

type Criticality string

/*
A (7x24)
A+ (SMS 7x24)
B (5x12)
C (never)
*/

const (
	Criticality7x24    Criticality = "A"
	CriticalitySMS7x24 Criticality = "A+"
	Criticality5x12    Criticality = "B"
	CriticalityNever   Criticality = "C"
)

func (c Criticality) Get() Criticality {
	if len(c) < 1 {
		return CriticalityNever
	}
	return c
}
