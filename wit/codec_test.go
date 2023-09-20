package wit

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/ydnar/wasm-tools-go/internal/callerfs"
)

var update = flag.Bool("update", false, "update golden files")

func TestDecodeJSON(t *testing.T) {
	err := filepath.WalkDir(callerfs.Path("../testdata"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if !strings.HasSuffix(path, ".wit.json") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		res, err := DecodeJSON(f)
		if err != nil {
			return err
		}
		data := pretty.Sprint(res)
		compareOrWrite(t, path, data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func compareOrWrite(t *testing.T, path, data string) {
	golden := path + ".golden"
	if *update {
		err := os.WriteFile(golden, []byte(data), 0644)
		if err != nil {
			t.Error(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Error(err)
	}
	if string(want) != data {
		dmp := diffmatchpatch.New()
		dmp.PatchMargin = 3
		diffs := dmp.DiffMain(string(want), data, false)
		t.Errorf("value for %s did not match golden value %s:\n%v", path, golden, dmp.DiffPrettyText(diffs))
		// fmt.Fprintln(os.Stderr, dmp.DiffPrettyText(diffs))
	}
}
