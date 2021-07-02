package global

import (
	"ela/foundation/path"
	"os"
)

func Initialize() {
	os.MkdirAll(path.GetSystemApp(), 0740)
	os.MkdirAll(path.GetSystemWWW(), 0740)
	os.MkdirAll(path.GetDefaultBackupPath(), 0744)
	os.MkdirAll(path.GetSystemAppData(""), 0774)
	if path.HasExternal() {
		os.MkdirAll(path.GetExternalApp(), 0776)
		os.MkdirAll(path.GetExternalWWW(), 0776)
		os.MkdirAll(path.GetExternalAppData(""), 0774)
	}
}
