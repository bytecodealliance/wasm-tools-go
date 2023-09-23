package wit

import (
	"flag"
	"os"
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldenFiles(t *testing.T) {
	p := pp.New()
	p.SetExportedOnly(true)
	p.SetColoringEnabled(false)

	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(path, func(t *testing.T) {
			data := p.Sprint(res)
			compareOrWrite(t, path, data)
		})
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

func compareOrWrite(t *testing.T, path, data string) {
	golden := path + ".golden"
	if *update {
		err := os.WriteFile(golden, []byte(data), 0644)
		if err != nil {
			t.Error(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Error(err)
		return
	}
	if string(want) != data {
		dmp := diffmatchpatch.New()
		dmp.PatchMargin = 3
		diffs := dmp.DiffMain(string(want), data, false)
		t.Errorf("value for %s did not match golden value %s:\n%v", path, golden, dmp.DiffPrettyText(diffs))
	}
}
