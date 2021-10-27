// this file handle any system components upgrading. Components can be outdated apps, registry, libraries etc

package sysupgrade

import "github.com/cansulting/elabox-system-tools/foundation/logger"

// starts upgrading the system.
// Components can be outdated apps, registry, libraries etc
func Start(oldBuild int, newBuild int) error {
	if oldBuild < 0 {
		logger.GetInstance().Debug().Msg("Fresh OS. System upgrade skipped.")
		return nil
	}
	logger.GetInstance().Debug().Msg("Starts upgrading system components...")
	upgrader := build3upgrade{}
	if err := upgrader.onUpgrade(oldBuild); err != nil {
		return err
	}
	logger.GetInstance().Debug().Msg("Component upgrades finished.")
	return nil
}
