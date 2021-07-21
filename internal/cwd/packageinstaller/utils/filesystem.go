package utils

import (
	"os"
	"strings"
)

// create directories if not exist
func ResolveDir(path string, perm os.FileMode) error {
	index := strings.LastIndex(path, "/")
	if index >= 0 {
		extractedDir := path[:index]
		err := os.MkdirAll(extractedDir, perm)
		return err
	}
	return nil
}
