package main

import (
	"testing"
)

func TestSilentInstall(test *testing.T) {
	//backup := Backup{}
	//error := backup.LoadAndApply("system.backup")
	//print(error.Error())
	//return
	newInstall := installer{backupEnabled: true}
	newInstall.initInstall(true)
	err := newInstall.decompress("../../builds/packages/system.ela")
	if err != nil {
		test.Error(err)
	}

}
