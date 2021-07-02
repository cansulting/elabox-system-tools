package main

import (
	"os"
	"strings"
)

// use filter to customize the destination or change name
type filter struct {
	keyword     string // for query
	replace     string // replace the path based from keyword. nil if not replace
	destination string // where file will be saved.
	additive    bool   // true if this filter will be mixed to other filter
}

// use to apply filter to path.
// resolveDir param is true if will create directories for new path
// return newpath, error, true if filter was applied
func (f *filter) applyTo(path string, perm os.FileMode) (string, error, bool) {
	contains := strings.Contains(path, f.keyword)
	if contains {
		var result string = path
		if f.replace != "" {
			result = strings.Replace(path, f.keyword, f.replace, 1)
		}
		if f.destination != "" {
			result = f.destination + "/" + result
		}
		var err error
		err = resolveDir(result, perm)
		return result, err, err == nil
	} else {
		return path, nil, false
	}
}
