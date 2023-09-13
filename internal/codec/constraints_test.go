package codec

import (
	"fmt"
	"testing"
)

func TestNumericTypeName(t *testing.T) {
	tests := []struct {
		v    fmt.Stringer
		want string
	}{
		{Numeric[int8]{}, "int8"},
		{Numeric[uint8]{}, "uint8"},
		{Numeric[int16]{}, "int16"},
		{Numeric[uint16]{}, "uint16"},
		{Numeric[int32]{}, "int32"},
		{Numeric[uint32]{}, "uint32"},
		{Numeric[int64]{}, "int64"},
		{Numeric[uint64]{}, "uint64"},
		{Numeric[float32]{}, "float32"},
		{Numeric[float64]{}, "float64"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.v.String()
			if got != tt.want {
				t.Errorf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
