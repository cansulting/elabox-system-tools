package app

import (
	"ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/registry/util"
	"log"
)

/*
	app.go
	Methods that provides access to application retrieval, registration etc
*/

// retrieve all packages
func RetrieveAllPackages(id string) ([]*data.PackageConfig, error) {
	row, err := retrievePackagesRaw(id, []string{"id, source, name, location"})
	if err != nil {
		return nil, err
	}
	defer row.Close()
	results := make([]*data.PackageConfig, 0, 10)
	if row.Next() {
		pk := data.DefaultPackage()
		row.Scan(&pk.PackageId, &pk.Source, &pk.Name, &pk.InstallLocation)
		results = append(results, pk)
	}
	return results, nil
}

// add package data to db
func RegisterPackage(pkData *data.PackageConfig) error {
	log.Println("Registering package " + pkData.PackageId)
	query := `
		replace into 
		packages(id, location, build, version, name, desc, source) 
		values(?,?,?,?,?,?,?)`
	err := util.ExecuteQuery(
		query,
		pkData.PackageId,
		pkData.InstallLocation,
		pkData.Build,
		pkData.Version,
		pkData.Name,
		pkData.Description,
		pkData.Source,
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
	log.Println("appman.RegisterPackageSrc success!")
	return config, nil
}

func RetrievePackage(id string) (*data.PackageConfig, error) {
	pks, err := RetrieveAllPackages(id)
	if err != nil {
		return nil, errors.SystemNew("appman.RetrievePackage failed", err)
	}
	if len(pks) > 0 {
		return pks[0], nil
	}
	return nil, nil
}

func RetrievePackagesWithActivity(action string) ([]string, error) {
	return RetrievePackagesForActivity(action)
}
