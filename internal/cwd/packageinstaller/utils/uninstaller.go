package utils

import (
	"ela/foundation/errors"
	"ela/foundation/logger"
	"ela/registry/app"
	"os"
)

// delete package based package id
func UninstallPackage(packageId string, deleteData bool) error {
	logger.GetInstance().Debug().Msg("Deleting old package " + packageId)
	// step: retrieve package location
	pk, err := app.RetrievePackage(packageId)
	if err != nil {
		return errors.SystemNew("Delete package failed "+packageId, err)
	}
	// step: not yet installed? skip
	if pk == nil {
		return nil
	}
	location := pk.GetInstallDir()
	// is file not exist. skip
	if _, err := os.Stat(location); err != nil {
		//log.Println("")
		return nil
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
	return nil
}
