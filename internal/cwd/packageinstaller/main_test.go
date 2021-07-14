package main

import (
	"ela/foundation/app"
	appd "ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/event/data"
	"ela/foundation/path"
	testapp "ela/foundation/testing/app"
	"ela/internal/cwd/packageinstaller/global"
	reg "ela/registry/app"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const pkid = "ela.installer"

// use to build installer. and put it to installation path
func TestBuildInstaller(test *testing.T) {
	log.Println("Building installer")
	buildFname := "packageinstaller.exe"
	//sourcePath := "./internal/cwd/packageinstaller"
	outputPath := "../../builds"
	targetPath := path.GetAppMain(pkid, false)
	cmd := exec.Command("go", "install")
	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		test.Error(err)
		return
	}
	if err := os.Rename(outputPath+"/bins/"+buildFname, targetPath); err != nil {
		test.Error(err)
	}
}

// step: create a copy of this application
func copyInstallerBinary(pkconfig *appd.PackageConfig) (string, error) {
	binPath := pkconfig.GetMainExec()
	dest := path.GetCacheDir() + "/" + filepath.Base(binPath)
	log.Println("Cloning binary @ " + dest)
	bytes, err := ioutil.ReadFile(binPath)
	if err != nil {
		return "", errors.SystemNew("installer.startCommandlineMode Failed to copy installer binary.", err)
	}
	if err := ioutil.WriteFile(dest, bytes, 0770); err != nil {
		return "", errors.SystemNew("installer.startCommandlineMode Failed to copy installer binary.", err)
	}
	return dest, nil
}

// test install a package and register it
func TestSystemUpdateCommandline(test *testing.T) {
	wd, _ := os.Getwd()
	pkpath := wd + "/../../builds/windows/packager/ela.installer.ela"
	pki, err := reg.RetrievePackage(pkid)
	if err != nil {
		test.Error(err)
		return
	}
	// if err := reg.CloseDB(); err != nil {
	// 	test.Error(err)
	// 	return
	// }
	dest, err := copyInstallerBinary(pki)
	if err != nil {
		log.Println(err.Error())
		test.Error(err)
		return
	}
	cmd := exec.Command(dest, pkpath)
	cmd.Dir = filepath.Dir(dest)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		test.Error(err)
		return
	}
	log.Println(string(bytes))
}

// test installer via activity
func TestRunActivityManually(t *testing.T) {
	InitializePath()
	pkgPath := "../../builds/windows/packager/ela.system.ela"
	//pkgPath := `C:\Users\Jhoemar\Documents\Projects\Elabox\system-tools\internal\builds\packages`
	controller, err := app.NewController(&activity{}, nil)
	if err != nil {
		t.Error(err)
	}
	global.AppController = controller
	testapp.RunTestApp(
		controller,
		data.ActionGroup{
			Activity: data.NewAction(constants.ACTION_APP_INSTALLER, "", pkgPath)})
}
