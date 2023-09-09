package wit

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ydnar/wit-bindgen-go/internal/testutil"
)

func TestDecodeJSON(t *testing.T) {
	f, err := os.Open(testutil.Path("../testdata/functions.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res, err := DecodeJSON(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", res)

	j, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(j))
}
