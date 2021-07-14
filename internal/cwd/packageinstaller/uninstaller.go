package main

import (
	"ela/foundation/errors"
	"ela/registry/app"
	"log"
	"os"
)

// delete package based from installed information
func deletePackage(packageId string, deleteData bool) error {
	log.Println("Deleting old package", packageId)
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
		log.Println("")
		return nil
	}
	// step: remove directory
	if err := os.RemoveAll(location); err != nil {
		return errors.SystemNew("Delete package failed "+packageId, err)
	}
	// step: remove data
	if deleteData {
		if err := os.RemoveAll(pk.GetDataDir()); err != nil {
			log.Println("Remove pacakage data failed", err.Error())
			return nil
		}
	}
	return nil
}
