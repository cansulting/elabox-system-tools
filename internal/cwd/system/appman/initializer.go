package appman

import (
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
)

func Initialize(commandline bool) error {
	if !commandline {
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
