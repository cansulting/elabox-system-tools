// provides utils for handling package history

package sysinstall

import (
	"io"
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
)

var OLD_PK = path.GetSystemAppDirData(constants.SYSTEM_SERVICE_ID) + "/" + "old_info.json"
var CUR_PK = path.GetSystemAppDir() + "/" + constants.SYSTEM_SERVICE_ID + "/" + constants.APP_CONFIG_NAME

func GetInstalledPackage() *data.PackageConfig {
	pk := data.DefaultPackage()
	pk.LoadFromSrc(CUR_PK)
	return pk
}

func GetOldPackage() *data.PackageConfig {
	pk := data.DefaultPackage()
	pk.LoadFromSrc(OLD_PK)
	return pk
}

func HasOldPackage() bool {
	_, err := os.Stat(OLD_PK)
	return err == nil
}

// use to create old package
func CreateOldPackageInfo() error {
	wfile, err := os.OpenFile(OLD_PK, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm.PUBLIC_VIEW)
	if err != nil {
		return err
	}
	rf, err := os.OpenFile(CUR_PK, os.O_RDONLY, perm.PUBLIC_VIEW)
	if err != nil {
		return err
	}
	if _, err2 := io.Copy(wfile, rf); err2 != nil {
		return err2
	}
	wfile.Sync()
	wfile.Close()
	rf.Close()
	return nil
}
