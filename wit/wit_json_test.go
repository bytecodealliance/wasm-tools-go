package wit

import (
	"encoding/json"
	"testing"
)

func TestJSONUnmarshalSizing(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		worlds     int
		interfaces int
		typedefs   int
		packages   int
	}{
		{"no data", `{}`, 0, 0, 0, 0},
		{"1 of each", `{"worlds": [{}], "interfaces": [{}], "types": [{}], "packages": [{}]}`, 1, 1, 1, 1},
		{"5 types", `{"worlds": [{}], "interfaces": [{}], "types": [{}, {}, {}, {}, {}], "packages": [{}]}`, 1, 1, 5, 1},
		{"2 interfaces, 2 packages", `{"worlds": [{}], "interfaces": [{}, {}], "types": [{}, {}, {}, {}, {}], "packages": [{}, {}]}`, 1, 2, 5, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res Resolve
			err := json.Unmarshal([]byte(tt.data), &res)
			if err != nil {
				t.Fatal(err)
			}
			if len(res.Worlds) != tt.worlds {
				t.Errorf("expected %d worlds, got %d", tt.worlds, len(res.Worlds))
			}
			if len(res.Interfaces) != tt.interfaces {
				t.Errorf("expected %d interfaces, got %d", tt.interfaces, len(res.Interfaces))
			}
			if len(res.TypeDefs) != tt.typedefs {
				t.Errorf("expected %d types, got %d", tt.typedefs, len(res.TypeDefs))
			}
			if len(res.Packages) != tt.packages {
				t.Errorf("expected %d packages, got %d", tt.packages, len(res.Packages))
			}
		})
	}
}
