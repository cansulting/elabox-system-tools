package system_update

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

// step: create a copy of this application
func CopyInstallerBinary(pkconfig *data.PackageConfig) (string, error) {
	binPath := path.GetAppMain(pkconfig.PackageId, !pkconfig.IsSystemPackage())
	dest := path.GetCacheDir() + "/" + filepath.Base(binPath)
	log.Println("Cloning binary @ " + dest)
	bytes, err := ioutil.ReadFile(binPath)
	if err != nil {
		return "", errors.SystemNew("installer.startCommandlineMode Failed to read installer binary.", err)
	}
	if err := ioutil.WriteFile(dest, bytes, 0770); err != nil {
		return "", errors.SystemNew("installer.startCommandlineMode Failed to write installer binary.", err)
	}
	return dest, nil
}

// install via commandline mode. this install the package safely by running process outside the system.
// @pkgPath is the location of package to be installed
// @pkg is package of info of installer
func Start(pkgPath string, pkg *data.PackageConfig) error {
	log.Println("Starting system update")
	// step: create a copy of this application and start exec
	dest, err := CopyInstallerBinary(pkg)
	if err != nil {
		return err
	}
	// step: run the copied binary
	cmd := exec.Command(dest, pkgPath)
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	println("ERROR " + err.Error())
	//}
	//println(string(out))

	if err := cmd.Start(); err != nil {
		return errors.SystemNew("installer.startCommandlineMode Failed to execute commandline binary.", err)
	}
	log.Println("Running commandline installer")
	time.Sleep(2 * time.Second)
	os.Exit(0)
	return nil
}
