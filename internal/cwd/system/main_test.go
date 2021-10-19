package main

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/cansulting/elabox-system-tools/internal/cwd/system/system_update"
	reg "github.com/cansulting/elabox-system-tools/registry/app"
)

// test install a package and register it
func TestCommandline(test *testing.T) {
	wd, _ := os.Getwd()
	pkpath := wd + "/../../builds/packages/packageinstaller.ela"
	pki, err := reg.RetrievePackage("ela.installer")
	if err != nil {
		test.Error(err)
		return
	}
	// if err := reg.CloseDB(); err != nil {
	// 	test.Error(err)
	// 	return
	// }
	dest, err := system_update.CopyInstallerBinary(pki)
	if err != nil {
		log.Println(err.Error())
		test.Error(err)
		return
	}
	cmd := exec.Command(dest, pkpath)
	cmd.Dir = pki.GetInstallDir()
	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		test.Error(err)
		return
	}

	/*
		newInstall := installer{BackupEnabled: true, SilentInstall: true}
		err := newInstall.Decompress("../../builds/packages/packageinstaller.ela")
		if err != nil {
			test.Error(err)
		}
		if err := newInstall.RegisterPackage(); err != nil {
			test.Error(err)
			return
		}*/
}
