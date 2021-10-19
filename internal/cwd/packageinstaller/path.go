package main

import (
	"log"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
)

func InitializePath() {
	if err := os.MkdirAll(path.GetSystemAppDir(), perm.PRIVATE); err != nil {
		log.Fatalln("Unable to create directories", err)
	}
	os.MkdirAll(path.GetSystemWWW(), perm.PRIVATE)
	os.MkdirAll(path.GetDefaultBackupPath(), perm.PUBLIC)
	os.MkdirAll(path.GetSystemAppDirData(""), perm.PUBLIC_WRITE)
	os.MkdirAll(path.GetCacheDir(), perm.PUBLIC)
	os.MkdirAll(path.GetLibPath(), perm.PUBLIC_VIEW)
	if path.HasExternal() {
		os.MkdirAll(path.GetExternalAppDir(), perm.PRIVATE)
		os.MkdirAll(path.GetExternalWWW(), perm.PRIVATE)
		os.MkdirAll(path.GetExternalAppData(""), perm.PUBLIC_WRITE)
	}
}
