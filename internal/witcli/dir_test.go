package witcli

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFindOrCreateDirFind(t *testing.T) {
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

	// Skip remaining tests on incompatible platforms
	if runtime.Compiler == "tinygo" || strings.Contains(runtime.GOARCH, "wasm") {
		return
	}

	gotPerm := info.Mode().Perm()
	wantPerm := fi.Mode().Perm()
	if gotPerm != wantPerm {
		t.Errorf("FindOrCreateDirExisting returned directory with file permissions %v; want %v", gotPerm, wantPerm)
	}
}

func TestFindOrCreateDirCreate(t *testing.T) {
	temp := t.TempDir()

	dir := filepath.Join(temp, "new")

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist, but it did.", dir)
	}

	info, err := FindOrCreateDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Skip remaining tests on incompatible platforms
	if runtime.Compiler == "tinygo" || strings.Contains(runtime.GOARCH, "wasm") {
		return
	}

	gotPerm := info.Mode().Perm()
	wantPerm := os.FileMode(0o755)
	if gotPerm != wantPerm {
		t.Errorf("FindOrCreateDirExisting returned directory with file permissions %v; want %v", gotPerm, wantPerm)
	}
}

func TestFindOrCreateDirPermissionDenied(t *testing.T) {
	// Skip test on incompatible platforms
	if runtime.Compiler == "tinygo" || strings.Contains(runtime.GOARCH, "wasm") {
		return
	}

	temp := t.TempDir()

	protected := filepath.Join(temp, "protected")
	if err := os.Mkdir(protected, os.FileMode(0o400)); err != nil {
		t.Fatal(err)
	}

	dir := filepath.Join(protected, "permission-denied")

	_, err := FindOrCreateDir(dir)
	if !errors.Is(err, fs.ErrPermission) {
		t.Errorf("Expected creation within protected directory to fail, but it did not: %v", err)
	}
}
