package main

import (
	"ela/foundation/path"
	"log"
	"os"
)

func InitializePath() {
	if err := os.MkdirAll(path.GetSystemApp(), 0740); err != nil {
		log.Fatalln("Unable to create directories", err)
	}
	os.MkdirAll(path.GetSystemWWW(), 0740)
	os.MkdirAll(path.GetDefaultBackupPath(), 0744)
	os.MkdirAll(path.GetSystemAppData(""), 0774)
	os.MkdirAll(path.GetCacheDir(), 0777)
	if path.HasExternal() {
		os.MkdirAll(path.GetExternalApp(), 0776)
		os.MkdirAll(path.GetExternalWWW(), 0776)
		os.MkdirAll(path.GetExternalAppData(""), 0774)
	}
}
