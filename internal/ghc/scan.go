package ghc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var errSrcIsNotDirectory = errors.New("source must be a directory")

func ScanDir(root string) ([]string, error) {
	if err := isValid(root); err != nil {
		return nil, err
	}

	directories := make([]string, 0)
	err := filepath.Walk(root, walkFnx(&directories))
	if err != nil {
		return nil, err
	}

	return directories, nil
}

func isValid(root string) (err error) {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errSrcIsNotDirectory
	}

	return nil
}

func walkFnx(directories *[]string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			*directories = append(*directories, path)
		}
		return nil
	}
}
