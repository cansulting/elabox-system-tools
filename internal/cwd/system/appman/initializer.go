package appman

import "github.com/cansulting/elabox-system-tools/internal/cwd/system/global"

func Initialize(commandline bool) error {
	if !commandline {
		if global.RUN_STARTUPAPPS {
			InitializeStartups()
		} else {
			global.Logger.Debug().Msg("Startup apps was disabled.")
		}
	}
	return nil
}
