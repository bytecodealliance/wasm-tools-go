package codec

import (
	"testing"
)

func TestTypeName(t *testing.T) {
	tests := []struct {
		f    func() string
		want string
	}{
		{TypeName[int8], "int8"},
		{TypeName[uint8], "uint8"},
		{TypeName[int16], "int16"},
		{TypeName[uint16], "uint16"},
		{TypeName[int32], "int32"},
		{TypeName[uint32], "uint32"},
		{TypeName[int64], "int64"},
		{TypeName[uint64], "uint64"},
		{TypeName[float32], "float32"},
		{TypeName[float64], "float64"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.f()
			if got != tt.want {
				t.Errorf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
