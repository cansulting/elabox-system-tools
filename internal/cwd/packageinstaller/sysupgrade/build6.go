// system upgrade for build 6 and lower
package sysupgrade

import (
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/env"
)

type build6 struct {
}

func (instance *build6) onUpgrade(oldbuild int) error {
	// 6 and lower  build should mark as already configured.
	// update the system variable name "config" to 1
	if oldbuild <= 6 {
		logger.GetInstance().Debug().Msg("Starts upgrading build lower or equal to 6")
		if err := env.Init(); err != nil {
			logger.GetInstance().Error().Err(err).Msg("failed to initialize environment, upgrade skip for 6")
			return nil
		}
		env.SetEnv("config", "1")
	}
	return nil
}
