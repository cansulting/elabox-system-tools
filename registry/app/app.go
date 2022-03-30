// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package app

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

/*
	app.go
	Methods that provides access to application retrieval, registration etc
*/

// retrieve all packages
func RetrieveAllPackages() ([]*data.PackageConfig, error) {
	row, err := retrievePackagesRaw("", []string{"id, source, version, name, location, nodejs, program, build"})
	if err != nil {
		return nil, err
	}
	defer row.Close()
	results := convertRawToPackageConfig(row)
	return results, nil
}

// add package data to db
func RegisterPackage(pkData *data.PackageConfig) error {
	logger.GetInstance().Info().Str("category", "registry").Msg("Registering package " + pkData.PackageId)
	query := `
		replace into 
		packages(id, location, build, version, name, desc, source, nodejs, exportService, program) 
		values(?,?,?,?,?,?,?,?,?,?)`
	err := util.ExecuteQuery(
		query,
		pkData.PackageId,
		pkData.InstallLocation,
		pkData.Build,
		pkData.Version,
		pkData.Name,
		pkData.Description,
		pkData.Source,
		pkData.Nodejs,
		pkData.ExportServices,
		pkData.Program,
	)
	if err != nil {
		return errors.SystemNew("records.AddPackage Failed to add "+pkData.PackageId, err)
	}
	if err := registerActivities(pkData); err != nil {
		return errors.SystemNew("appman.RegisterPackageSrc failed", err)
	}
	return nil
}

// use to register package base from dir location
func RegisterPackageSrc(srcDir string) (*data.PackageConfig, error) {
	config := data.DefaultPackage()
	if err := config.LoadFromSrc(srcDir + "/" + constants.APP_CONFIG_NAME); err != nil {
		return nil, errors.SystemNew("appman.RegisterPackageSrc couldnt load from source ", err)
	}
	if err := RegisterPackage(config); err != nil {
		return nil, errors.SystemNew("appman.RegisterPackageSrc failed", err)
	}
	//log.Println("appman.RegisterPackageSrc success!")
	return config, nil
}

// remove package data from db
func UnregisterPackage(pkId string) error {
	logger.GetInstance().Info().Str("category", "registry").Msg("Unregistering package " + pkId)
	query := `delete from packages where id = ?`
	err := util.ExecuteQuery(query, pkId)
	if err != nil {
		return errors.SystemNew("records.UnregisterPackage Failed to remove "+pkId, err)
	}
	if err := removeActivities(pkId); err != nil {
		return errors.SystemNew("appman.UnregisterPackage failed removing activities", err)
	}
	return nil
}

func RetrievePackage(id string) (*data.PackageConfig, error) {
	pks, err := retrievePackagesRaw(id, []string{"id", "source", "version", "name", "location", "nodejs", "program", "build"})
	if err != nil {
		return nil, errors.SystemNew("appman.RetrievePackage failed", err)
	}
	results := convertRawToPackageConfig(pks)
	if len(results) > 0 {
		return results[0], nil
	}
	return nil, nil
}

func RetrievePackagesWithActivity(action string) ([]string, error) {
	return RetrievePackagesForActivity(action)
}

// retrieve all packages that needs to execute upon startup
func RetrieveStartupPackages() ([]*data.PackageConfig, error) {
	row, err := retrievePackagesWhere(
		[]string{"id, source, name, location, nodejs, exportService, program"},
		"nodejs=true or exportService=true")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	results := make([]*data.PackageConfig, 0, 10)
	for row.Next() {
		pk := data.DefaultPackage()
		row.Scan(&pk.PackageId, &pk.Source, &pk.Name, &pk.InstallLocation, &pk.Nodejs, &pk.ExportServices, &pk.Program)
		results = append(results, pk)
	}
	return results, nil
}
