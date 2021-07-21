package utils

import (
	"io"
	"io/fs"
	"os"
	"strings"
)

// delegate type
type process func(string, io.ReadCloser, uint64) error

// use Filter files from package. use to move a file to specific location or make subprocess
type Filter struct {
	Keyword       string      // for query. add asteriest as prefix if keyword can be anywhere
	Rename        string      // use to replace the path based from keyword. nil if skip replace
	InstallTo     string      // where file will be saved.
	CustomProcess process     // to be called for custom processing for file
	Perm          fs.FileMode // the current file permission
	//additive    bool   // true if this Filter will be mixed to other Filter
}

// save the reader to specified path
func (instance *Filter) Save(path string, reader io.ReadCloser) error {
	// step: create dest file
	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, instance.Perm)
	if err != nil {
		return err
	}
	// step: write to file
	io.Copy(newFile, reader)
	newFile.Close()
	return nil
}

// use to apply Filter to path.
// return newpath, error, true if Filter was applied
func (f *Filter) CanApply(path string, reader io.ReadCloser, size uint64) (string, error, bool) {
	// path is absolute if contains not asterisk
	isAbsolute := true
	keyword := f.Keyword
	if keyword[0] == '*' {
		isAbsolute = false
		keyword = f.Keyword[1:]
	}
	indexFromPath := strings.Index(path, keyword)
	if indexFromPath >= 0 {
		// keyword is absolute but not found
		if isAbsolute && indexFromPath > 0 {
			return path, nil, false
		}
		// rename keyword
		var result string = path
		if f.Rename != "" {
			result = strings.Replace(path, f.Keyword, f.Rename, 1)
		}
		// step: execute custom process
		if f.CustomProcess != nil {
			f.CustomProcess(path, reader, size)
			return "", nil, true
		}
		// step: install to target location
		if f.InstallTo != "" {
			result = f.InstallTo + "/" + result
			err := ResolveDir(result, f.Perm)
			return result, err, err == nil
		}
		return "", nil, true
	} else {
		return path, nil, false
	}
}
