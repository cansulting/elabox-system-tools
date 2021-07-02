package main

import (
	"os"
	"strings"
)

func resolveDir(path string, perm os.FileMode) error {
	index := strings.LastIndex(path, "/")
	if index >= 0 {
		extractedDir := path[:index]
		err := os.MkdirAll(extractedDir, perm)
		return err
	}
	return nil
}
