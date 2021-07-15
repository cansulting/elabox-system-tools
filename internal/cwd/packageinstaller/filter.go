package main

import (
	"ela/internal/cwd/packageinstaller/utils"
	"io"
	"os"
	"strings"
)

// delegate type
type process func(string, io.ReadCloser, uint64) error

// use filter files from package. use to move a file to specific location or make subprocess
type filter struct {
	keyword       string  // for query. add asteriest as prefix if keyword can be anywhere
	rename        string  // use to replace the path based from keyword. nil if skip replace
	installTo     string  // where file will be saved.
	customProcess process // to be called for custom processing for file
	//additive    bool   // true if this filter will be mixed to other filter
}

// use to apply filter to path.
// return newpath, error, true if filter was applied
func (f *filter) applyTo(path string, perm os.FileMode, reader io.ReadCloser, size uint64) (string, error, bool) {
	// path is absolute if contains not asterisk
	isAbsolute := true
	keyword := f.keyword
	if keyword[0] == '*' {
		isAbsolute = false
		keyword = f.keyword[1:]
	}
	indexFromPath := strings.Index(path, keyword)
	if indexFromPath >= 0 {
		// keyword is absolute but not found
		if isAbsolute && indexFromPath > 0 {
			return path, nil, false
		}
		// rename keyword
		var result string = path
		if f.rename != "" {
			result = strings.Replace(path, f.keyword, f.rename, 1)
		}
		// step: execute custom process
		if f.customProcess != nil {
			f.customProcess(path, reader, size)
			return "", nil, true
		}
		// step: install to target location
		if f.installTo != "" {
			result = f.installTo + "/" + result
			err := utils.ResolveDir(result, perm)
			return result, err, err == nil
		}
		return "", nil, true
	} else {
		return path, nil, false
	}
}
