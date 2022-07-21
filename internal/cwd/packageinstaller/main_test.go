package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	appd "github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/path"

	//testapp "github.com/cansulting/elabox-system-tools/foundation/testing/app"

	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/pkg"
)

const pkid = "ela.installer"
const outputPath = "../../builds/linux/packageinstaller/bin/packageinstaller"

// use to build installer. and put it to installation path
func TestBuildInstaller(test *testing.T) {
	log.Println("Building installer")
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
	binPath := pkconfig.GetMainProgram()
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
	pkpath := wd + "/../../builds/linux/system/ela.system.box"
	/*pki, err := reg.RetrievePackage(pkid)
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
		}*/

	processInstallCommand(pkpath, false, true, false)
	//log.Println(string(bytes))
}

// test install a package and register it
func TestSystemUpdateCommandline2(test *testing.T) {
	wd, _ := os.Getwd()
	pkpath := wd + "/../../builds/linux/companion/ela.companion.box"
	pkg, err := pkg.LoadFromSource(pkpath)
	if err != nil {
		test.Error(err)
		return
	}
	newInstall := NewInstaller(pkg, true, false)
	// step: start install
	if err := newInstall.Start(); err != nil {
		test.Error(err)
		return
	}
}

//test uninstall app

func TestSystemUninstallCommandLine(test *testing.T) {
	err := processUninstallCommand("ela.sample", false)
	if err != nil {
		test.Error(err)
		return
	}
}

// test installer via activity
// func TestRunActivityManually(t *testing.T) {
// 	InitializePath()
// 	pkgPath := "../../builds/linux/system/ela.system.box"
// 	//pkgPath := `C:\Users\Jhoemar\Documents\Projects\Elabox\system-tools\internal\builds\packages`
// 	controller, err := app.NewController(&activity{}, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	pkc.AppController = controller
// 	action := data.NewAction(constants.ACTION_APP_SYSTEM_INSTALL, "", pkgPath)
// 	testapp.RunTestApp(
// 		controller,
// 		data.ActionGroup{
// 			Activity: &action,
// 		})
// }
