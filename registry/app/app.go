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
	row, err := retrievePackagesRaw("", []string{"id, source, version, name, location, nodejs, program"})
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
	if err := addActivities(pkData); err != nil {
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

func RetrievePackage(id string) (*data.PackageConfig, error) {
	pks, err := retrievePackagesRaw(id, []string{"id", "source", "version", "name", "location", "nodejs", "program"})
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
