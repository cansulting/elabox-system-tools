// this is upgrade for build number 3

package sysupgrade

import (
	"strconv"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

type build3upgrade struct {
}

func (instance build3upgrade) onUpgrade(oldBuild int) error {
	if oldBuild <= 3 {
		// we need to delete registry. registry was updated
		logger.GetInstance().Debug().Msg("Deleting app registry...")
		if err := util.DeleteDB(); err != nil {
			logger.GetInstance().Warn().Err(err).Msg("Failed deleting db.")
		}
		return nil
	}
	logger.GetInstance().Debug().Msg(
		"No upgrade found for old build " + strconv.Itoa(oldBuild) + ". skipped.")
	return nil
}
