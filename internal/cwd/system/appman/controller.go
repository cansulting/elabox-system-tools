package appman

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/records"
)

func Initialize(commandline bool) error {
	if !commandline {
		global.Connector.Subscribe(constants.SYSTEM_SERVICE_ID, OnRecievedRequest)
	}
	if err := records.Initialize(); err != nil {
		return errors.SystemNew("appman.Initialize failed to initialize records", err)
	}
	return nil
}
