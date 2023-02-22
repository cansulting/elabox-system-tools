package env

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/system"
)

const ELAENV = "ELAENV"         // environment variable for build mode
const ELAVERSION = "ELAVERSION" // current version of system

// intialize system configuration
func Init() error {
	if system.BuildMode != system.DEBUG {
		logger.ConsoleOut = false
	}
	logger.GetInstance().Info().Msg("System running " + string(system.BuildMode))
	if err := system.SetEnv(ELAENV, string(system.BuildMode)); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	pkg := data.DefaultPackage()
	err := pkg.LoadFromLocation(constants.SYSTEM_SERVICE_ID, data.SYSTEM)
	if pkg == nil {
		err = errors.SystemNew("unable to load package", nil)
	}
	if err != nil {
		return errors.SystemNew("System config environment error ", err)
	}
	if err := system.SetEnv(ELAVERSION, pkg.Version); err != nil {
		return errors.SystemNew("System Config Environment error", err)
	}
	return nil
}
