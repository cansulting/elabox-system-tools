package appman

import (
	"ela/foundation/constants"
	"ela/internal/cwd/system/global"
)

func Initialize(commandline bool) error {
	if !commandline {
		global.Connector.Subscribe(constants.SYSTEM_SERVICE_ID, OnRecievedRequest)
	}
	InitializeStartups()
	return nil
}
