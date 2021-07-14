package appman

import (
	system "ela/foundation/app/data"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/internal/cwd/system/records"
	"log"
)

/*
	library.go
	Methods that provides access to application retrieval, registration etc
*/

// use to register package base from dir location
func RegisterPackageSrc(srcDir string) (*system.PackageConfig, error) {
	config := system.DefaultPackage()
	if err := config.LoadFromSrc(srcDir + "/" + constants.APP_CONFIG_NAME); err != nil {
		return nil, errors.SystemNew("appman.RegisterPackageSrc couldnt load from source ", err)
	}
	if err := records.AddPackage(config); err != nil {
		return nil, errors.SystemNew("appman.RegisterPackageSrc failed", err)
	}
	if err := records.AddActivities(config); err != nil {
		return nil, errors.SystemNew("appman.RegisterPackageSrc failed", err)
	}
	log.Println("appman.RegisterPackageSrc success!")
	return config, nil
}

func RetrievePackage(id string) (*system.PackageConfig, error) {
	pks, err := records.RetrievePackages(id)
	if err != nil {
		return nil, errors.SystemNew("appman.RetrievePackage failed", err)
	}
	if len(pks) > 0 {
		return pks[0], nil
	}
	return nil, nil
}

func RetrieveAllPackages() ([]*system.PackageConfig, error) {
	return records.RetrievePackages("")
}

func RetrievePackagesWithActivity(action string) ([]string, error) {
	return records.RetrievePackagesForActivity(action)
}

/*
func RetrievePackagesWithBroadcast(action string) ([]string, error) {
	return records.RetrievePackagesWithBroadcast(action)
}*/

func retrievePackageSource(packageId string) (string, error) {
	pk, err := RetrievePackage(packageId)
	if err != nil {
		return "", err
	}
	if pk != nil {
		return pk.Source, nil
	}
	return "", nil
}
