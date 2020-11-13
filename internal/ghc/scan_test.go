package ghc_test

import (
	"fgh/internal/ghc"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func makeDir(t *testing.T, base, newDir string) {
	err := os.MkdirAll(filepath.Join(base, newDir), 0700)
	checkIfError(t, err)
}

func makeFile(t *testing.T, base, newFile string) {
	err := ioutil.WriteFile(filepath.Join(base, newFile), nil, 0600)
	checkIfError(t, err)
}

func TestDetectGitDirectories(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitDirectories")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	makeDir(t, dir, "a/.git")
	makeDir(t, dir, "b/.git")
	makeDir(t, dir, "c/not/relevant/directory")
	makeFile(t, dir, ".git")
	makeFile(t, dir, "c/not/relevant/directory/random_file")

	having, err := ghc.ScanDir(dir)
	checkIfError(t, err)

	expected := []string{
		filepath.Join(dir, "a/.git"),
		filepath.Join(dir, "b/.git"),
	}

	if !reflect.DeepEqual(expected, having) {
		t.Errorf("Expecting %v, having %v", expected, having)
	}
}
