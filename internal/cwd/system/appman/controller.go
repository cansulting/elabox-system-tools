package appman

import (
	"ela/foundation/constants"
	"ela/internal/cwd/system/records"
	"ela/internal/cwd/system/servicecenter"
)

func Initialize(commandline bool) error {
	servicecenter.RegisterService(constants.SYSTEM_SERVICE_ID, OnRecievedRequest)
	if err := records.Initialize(); err != nil {
		return err
	}
	return nil
}
