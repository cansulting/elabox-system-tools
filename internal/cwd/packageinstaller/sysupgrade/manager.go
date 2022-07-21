// this file handle any system components upgrading. Components can be outdated apps, registry, libraries etc

package sysupgrade

import (
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/sysinstall.go"
)

// starts upgrading the system.
// Components can be outdated apps, registry, libraries etc
func Start(oldBuild int, newBuild int) error {
	if oldBuild < 0 {
		logger.GetInstance().Debug().Msg("Fresh OS. System upgrade skipped.")
		return nil
	}
	upgrader := build3upgrade{}
	if err := upgrader.onUpgrade(oldBuild); err != nil {
		return err
	}
	logger.GetInstance().Debug().Msg("Component upgrades finished.")
	return nil
}

func CheckAndUpgrade(newBuildNum int) {
	// upgrades
	oldpk := sysinstall.GetInstalledPackage()
	oldbuildnum := -1
	if oldpk != nil {
		oldbuildnum = int(oldpk.Build)
	}
	if err := Start(oldbuildnum, newBuildNum); err != nil {
		constants.Logger.Error().Err(err).Stack().Caller().Msg("Failed system upgrade.")
		return
	}
}
