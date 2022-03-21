package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/errors"
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

// save the reader to specified path
func CopyToTarget(target string, reader io.ReadCloser, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(target), perm); err != nil {
		return errors.SystemNew("Failed copying to target.", err)
	}
	// step: create dest file
	newFile, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer newFile.Close()
	// step: write to file
	_, err = io.Copy(newFile, reader)
	return err
}
