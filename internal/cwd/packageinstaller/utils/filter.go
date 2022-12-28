package utils

import (
	"io"
	"io/fs"
	"os"
	"path"
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
	LinkTo        string      // directory for new link file
}

// save the reader to specified path
func (instance *Filter) Save(_path string, reader io.ReadCloser) error {
	if err := CopyToTarget(_path, reader, instance.Perm); err != nil {
		return err
	}
	// create link
	if instance.LinkTo != "" {
		link := instance.LinkTo + "/" + path.Base(_path)
		// remove any existing link file
		if _, err := os.Stat(link); err == nil {
			os.Remove(link)
		}
		return os.Symlink(_path, link)
	}
	return nil
}

// use to apply Filter to path.
// return newpath, error, true if Filter was applied
func (f *Filter) CanApply(ppath string, reader io.ReadCloser, size uint64) (string, error, bool) {
	// path is absolute if contains not asterisk
	isAbsolute := true
	keyword := f.Keyword
	if keyword[0] == '*' {
		isAbsolute = false
		keyword = f.Keyword[1:]
	}
	indexFromPath := strings.Index(ppath, keyword)
	// this keyword found on specified path?
	if indexFromPath >= 0 {
		// keyword is absolute but not found
		if isAbsolute && indexFromPath > 0 {
			return ppath, nil, false
		}
		// rename keyword
		var result string = ppath
		if f.Rename != "" {
			result = strings.Replace(ppath, f.Keyword, f.Rename, 1)
		}
		// step: execute custom process
		if f.CustomProcess != nil {
			f.CustomProcess(ppath, reader, size)
			return "", nil, true
		}
		// step: create symlink
		if f.LinkTo != "" {
			if err := ResolveDir(f.LinkTo+"/", f.Perm); err != nil {
				return result, err, true
			}
		}
		// step: install to target location
		if f.InstallTo != "" {
			result = f.InstallTo + "/" + result
			err := ResolveDir(result, f.Perm)
			return result, err, err == nil
		}
		return "", nil, true
	} else {
		return ppath, nil, false
	}
}
