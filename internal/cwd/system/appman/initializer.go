package appman

import "github.com/cansulting/elabox-system-tools/internal/cwd/system/global"

func Initialize(commandline bool) error {
	if !commandline && global.RUN_STARTUPAPPS {
		InitializeStartups()
	}
	return nil
}
