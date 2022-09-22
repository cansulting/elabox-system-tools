package env

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

const ELAENV = "ELAENV"         // environment variable for build mode
const ELAVERSION = "ELAVERSION" // current version of system
const ELASHUTDOWNSTATUS="ELASHUTDOWNSTATUS"

// intialize system configuration
func Init() error {
	if system.BuildMode != system.DEBUG {
		logger.ConsoleOut = false
	}
	logger.GetInstance().Info().Msg("System running " + string(system.BuildMode))
	if err := system.SetEnv(ELAENV, string(system.BuildMode)); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	pkg, err := app.RetrievePackage(constants.SYSTEM_SERVICE_ID)
	if pkg == nil {
		err = errors.SystemNew("unable to load package", nil)
	}
	if err != nil {
		return errors.SystemNew("System config environment error ", err)
	}
	if err := system.SetEnv(ELAVERSION, pkg.Version); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	if(system.GetEnv(ELASHUTDOWNSTATUS) != "properly_shutdown") {
		if err := system.SetEnv(ELASHUTDOWNSTATUS, "not_properly_shutdown"); err != nil {
			return errors.SystemNew("System Config Environment error", err)
		}	
	}
	return nil
}
