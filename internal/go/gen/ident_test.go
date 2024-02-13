package gen

import "testing"

func TestParseSelector(t *testing.T) {
	tests := []struct {
		s    string
		path string
		name string
	}{
		{"io", "io", "io"},
		{"io/fs", "io/fs", "fs"},
		{"encoding/json", "encoding/json", "json"},
		{"encoding/xml", "encoding/xml", "xml"},
		{"encoding/xml#xml", "encoding/xml", "xml"},
		{"encoding/xml#encxml", "encoding/xml", "encxml"},
		{"encoding/xml#Encoder", "encoding/xml", "Encoder"},
		{"encoding/xml#Decoder", "encoding/xml", "Decoder"},
		{"wasi/clocks", "wasi/clocks", "clocks"},
		{"wasi/clocks#clocks", "wasi/clocks", "clocks"},
		{"wasi/clocks/wallclock", "wasi/clocks/wallclock", "wallclock"},
		{"wasi/clocks/wallclock#wallclock", "wasi/clocks/wallclock", "wallclock"},
		{"wasi/clocks/wallclock#DateTime", "wasi/clocks/wallclock", "DateTime"},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			path, name := ParseSelector(tt.s)
			if path != tt.path || name != tt.name {
				t.Errorf("ParseSelector(%q): %q, %q; expected %q, %q", tt.s, path, name, tt.path, tt.name)
			}
		})
	}
}
