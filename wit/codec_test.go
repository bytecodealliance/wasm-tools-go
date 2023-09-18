package wit

import (
	"fmt"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/ydnar/wasm-tools-go/internal/testutil"
)

func TestDecodeJSON(t *testing.T) {
	f, err := os.Open(testutil.Path("../testdata/worlds-with-types.wit.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res, err := DecodeJSON(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%# v\n", pretty.Formatter(res))
}
