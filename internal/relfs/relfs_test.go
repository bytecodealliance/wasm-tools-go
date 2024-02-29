package relfs

import (
	"path/filepath"
	"testing"
)

func TestGetwd(t *testing.T) {
	cwd, err := Getwd()
	if err != nil {
		t.Error(err)
	}
	t.Logf("WD: %s", cwd)
	if got, want := filepath.Base(cwd), "relfs"; got != want {
		t.Errorf("WD: got base %s, expected %s", got, want)
	}
}
