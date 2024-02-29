package relpath

import (
	"path/filepath"
	"testing"
)

func TestAbs(t *testing.T) {
	path, err := Abs(".")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Abs: %s", path)
	if got, want := filepath.Base(path), "relpath"; got != want {
		t.Errorf("Abs: got base %s, expected %s", got, want)
	}

	path, err = Abs("..")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Abs: %s", path)
	if got, want := filepath.Base(path), "internal"; got != want {
		t.Errorf("Abs: got base %s, expected %s", got, want)
	}
}

func TestGetwd(t *testing.T) {
	path, err := Getwd()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Getwd: %s", path)
	if got, want := filepath.Base(path), "relpath"; got != want {
		t.Errorf("Getwd: got base %s, expected %s", got, want)
	}
}
