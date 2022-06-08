package utils

import (
	"os"

	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/broadcast"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

// delete package based package id
func UninstallPackage(
	packageId string,
	deleteData bool,
	unregister bool,
	broadcastUpdate bool) error {
	logger.GetInstance().Debug().Msg("Deleting old package " + packageId)
	// step: retrieve package location
	pk, err := app.RetrievePackage(packageId)
	if err != nil {
		return errors.SystemNew("Delete package failed "+packageId, err)
	}
	// step: not yet installed? skip
	if pk == nil {
		logger.GetInstance().Debug().Msg(packageId + " package already removed skipping")
		return nil
	}
	location := pk.GetInstallDir()
	// is file not exist. skip
	if _, err := os.Stat(location); err != nil {
		//log.Println("")
		//return nil
	}
	// step: remove app directory
	if err := os.RemoveAll(location); err != nil {
		return errors.SystemNew("Delete package failed "+packageId, err)
	}
	// step: remove app data
	if deleteData {
		if err := os.RemoveAll(pk.GetDataDir()); err != nil {
			logger.GetInstance().Error().Err(err).Caller().Msg("Remove pacakage data failed")
			return nil
		}
	}
	// step: remove www dir
	www := pk.GetWWWDir()
	if _, err := os.Stat(www); err == nil {
		if err := os.RemoveAll(www); err != nil {
			return errors.SystemNew("failed to delete www dir for "+packageId, err)
		}
	}
	if unregister {
		// step: unregister package
		if err := app.UnregisterPackage(packageId); err != nil {
			return errors.SystemNew("failed to unregister package "+packageId, err)
		}
	}
	if broadcastUpdate {
		// step: broadcast update
		broadcast.UpdateSystem(packageId, broadcast.UNINSTALLED)
	}
	return nil
}
