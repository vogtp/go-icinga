package check

import (
	"testing"
)

func Test_parseFloat(t *testing.T) {
	var f32 float32 = 1.5
	var i16 int16 = 5
	var i32 int32 = 5
	var i64 int64 = 5
	tests := []struct {
		name string
		args any
		want float64
		ok   bool
	}{
		{name: "string 5", args: "5", want: 5, ok: true},
		{name: "string 5%", args: "5%", want: 0, ok: false},
		{name: "float64", args: 1.5, want: 1.5, ok: true},
		{name: "float32", args: f32, want: 1.5, ok: true},
		{name: "int", args: 5, want: 5, ok: true},
		{name: "int16", args: i16, want: 5, ok: true},
		{name: "int32", args: i32, want: 5, ok: true},
		{name: "int64", args: i64, want: 5, ok: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseFloat(tt.args)
			if got != tt.want {
				t.Errorf("parseFloat() float = %v, want %v", got, tt.want)
			}
			if got1 != tt.ok {
				t.Errorf("parseFloat() ok = %v, want %v", got1, tt.ok)
			}
		})
	}
}
