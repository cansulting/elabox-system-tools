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
	outputPath := "../../builds/linux/bins/packageinstaller"
	cmd := exec.Command("go", "build", "-o", outputPath)
	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		test.Error(err)
		return
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
	pkpath := wd + "/../../builds/linux/packager/ela.system.box"
	pki, err := reg.RetrievePackage(pkid)
	if err != nil {
		log.Println(err)
		test.Error(err)
		return
	}
	dest, err := copyInstallerBinary(pki)
	if err != nil {
		log.Println(err.Error())
		test.Error(err)
		return
	}
	cmd := exec.Command("sudo "+dest, pkpath)
	cmd.Dir = filepath.Dir(dest)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		test.Error(err)
		return
	}
	log.Println(string(bytes))
}

// test install a package and register it
func TestSystemUpdateCommandline2(test *testing.T) {
	wd, _ := os.Getwd()
	pkpath := wd + "/../../builds/linux/packager/ela.system.box"
	newInstall := installer{BackupEnabled: true, RunCustomInstaller: true}
	// step: start install
	if err := newInstall.Decompress(pkpath); err != nil {
		log.Fatal(err.Error())
		return
	}
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
