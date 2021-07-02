package main

import (
	"testing"
)

// test install a package and register it
func TestSilentInstall(test *testing.T) {
	//backup := Backup{}
	//error := backup.LoadAndApply("system.backup")
	//print(error.Error())
	//return
	newInstall := installer{BackupEnabled: true, SilentInstall: true}
	err := newInstall.decompress("../../builds/packages/system.ela")
	if err != nil {
		test.Error(err)
	}
	if err := newInstall.registerPackage(); err != nil {
		test.Error(err)
		return
	}
}
