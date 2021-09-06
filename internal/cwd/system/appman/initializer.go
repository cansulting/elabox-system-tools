package appman

import "ela/internal/cwd/system/global"

func Initialize(commandline bool) error {
	if !commandline && global.RUN_STARTUPAPPS {
		InitializeStartups()
	}
	return nil
}
