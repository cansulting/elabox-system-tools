package main

import (
	"ela/foundation/path"
	"ela/foundation/perm"
	"log"
	"os"
)

func InitializePath() {
	if err := os.MkdirAll(path.GetSystemAppDir(), perm.PRIVATE); err != nil {
		log.Fatalln("Unable to create directories", err)
	}
	os.MkdirAll(path.GetSystemWWW(), perm.PRIVATE)
	os.MkdirAll(path.GetDefaultBackupPath(), perm.PUBLIC_VIEW)
	os.MkdirAll(path.GetSystemAppDirData(""), perm.PUBLIC_WRITE)
	os.MkdirAll(path.GetCacheDir(), perm.PUBLIC)
	os.MkdirAll(path.GetLibPath(), perm.PUBLIC_VIEW)
	if path.HasExternal() {
		os.MkdirAll(path.GetExternalAppDir(), perm.PRIVATE)
		os.MkdirAll(path.GetExternalWWW(), perm.PRIVATE)
		os.MkdirAll(path.GetExternalAppData(""), perm.PUBLIC_WRITE)
	}
}
