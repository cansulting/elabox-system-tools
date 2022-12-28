package appman

import (
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
)

func Initialize(commandline bool) error {
	if !commandline {
		initDirectories()
		if global.RUN_STARTUPAPPS {
			// if system is not yet configured only run utility services needed to run the basic functions
			if system.GetEnv(global.CONFIG_ENV) == "1" {
				InitializeAllPackages()
			} else {
				return initializeConfigPackages()
			}
		} else {
			global.Logger.Debug().Msg("Startup apps was disabled.")
		}
	}
	return nil
}

func initDirectories() {
	os.MkdirAll(path.PATH_USERS, perm.PUBLIC)
	os.MkdirAll(path.PATH_HOME_APPS, perm.PUBLIC_VIEW)
	os.MkdirAll(path.PATH_HOME_DATA, perm.PRIVATE)
	os.MkdirAll(path.PATH_HOME_DOCUMENTS, perm.PUBLIC)
}
