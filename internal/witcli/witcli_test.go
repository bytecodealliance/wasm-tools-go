//go:build !tinygo

package witcli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindOrCreateDir_Find(t *testing.T) {
	temp := t.TempDir()

	dir := filepath.Join(temp, "existing")
	if err := os.Mkdir(dir, os.ModeDir|os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := FindOrCreateDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	gotModTime := info.ModTime()
	wantModTime := fi.ModTime()
	if gotModTime != wantModTime {
		t.Errorf("FindOrCreateDirExisting returned directory with modtime %v; want %v", gotModTime, wantModTime)
	}

	gotPerm := info.Mode().Perm()
	wantPerm := fi.Mode().Perm()
	if gotPerm != wantPerm {
		t.Errorf("FindOrCreateDirExisting returned directory with file permissions %v; want %v", gotPerm, wantPerm)
	}
}

func TestFindOrCreateDir_Create(t *testing.T) {
	temp := t.TempDir()

	dir := filepath.Join(temp, "new")

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist, but it did.", dir)
	}

	info, err := FindOrCreateDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	gotPerm := int(info.Mode().Perm())
	wantPerm := 0o755
	if gotPerm != wantPerm {
		t.Errorf("FindOrCreateDirExisting returned directory with file permissions %v; want %v", gotPerm, wantPerm)
	}
}
