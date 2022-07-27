package config

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

const ELAENV = "ELAENV"         // environment variable for build mode
const ELAVERSION = "ELAVERSION" // current version of system

// intialize system configuration
func Init() error {
	if GetBuildMode() != DEBUG {
		logger.ConsoleOut = false
	}

	if err := SetEnv(ELAENV, string(GetBuildMode())); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	pkg, err := app.RetrievePackage(constants.SYSTEM_SERVICE_ID)
	if err != nil {
		return errors.SystemNew("System config environment error ", err)
	}
	if err := SetEnv(ELAVERSION, pkg.Version); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	return nil
}