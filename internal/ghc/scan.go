package ghc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func isValid(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("source must be a directory")
	}

	return nil
}

func ScanDir(root string) ([]string, error) {
	if err := isValid(root); err != nil {
		return nil, err
	}

	directories := make([]string, 0)
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			directories = append(directories, path)
		}
		return nil
	}

	err := filepath.Walk(root, walk)
	if err != nil {
		return nil, err
	}

	return directories, nil
}
