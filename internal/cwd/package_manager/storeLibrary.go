package main

import (
	"dashboard/package_manager/data"
	"dashboard/package_manager/services/installer"
	"errors"
	"sort"
	data2 "store/data"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

var deviceSerial = ""

// retrieve all apps
// @beta is true if include all apps for testing and demo apps
func RetrieveAllApps(beta bool) ([]data.PackageInfo, error) {
	// step: retrieve all apps from registry
	storeItems, err := app.
	if err != nil {
		return nil, errors.New("unable to retrieve all installed packages. inner: " + err.Error())
	}
	var previews = make([]data.PackageInfo, 0, len(storeItems))
	var tmpPreview data.PackageInfo
	// step: iterate on packages
	for _, pkg := range storeItems {
		installedInfo, err := app.RetrievePackage(pkg.Id)
		if err != nil {
			logger.GetInstance().Debug().Msg("unable to retrieve cache item for package: " + pkg.Id + ". inner: " + err.Error())
		}
		// if false && beta && pkg.Release.ReleaseType == data2.Beta && len(pkg.Tester.Users) > 0 {
		// 	tester, err := isTester(pkg.Tester.Users)
		// 	if err != nil {
		// 		logger.GetInstance().Debug().Msg("unable to validate user: " + err.Error())
		// 		continue
		// 	}
		// 	// not tester and not installed, skip the package
		// 	if !tester && installedInfo == nil {
		// 		continue
		// 	}
		// }
		tmpPreview = data.PackageInfo{}
		tmpPreview.AddInfo(installedInfo, &pkg, false)
		if task := installer.GetTask(tmpPreview.Id); task != nil {
			tmpPreview.Status = task.Status
		}

		// check if currently in download
		previews = append(previews, tmpPreview)
	}
	// sort preview by name
	sort.Slice(previews, func(i, j int) bool {
		return previews[i].Name < previews[j].Name
	})

	return previews, nil
}

// retrieve detailed information about an app
func RetrieveApp(pkgId string, storehubId string) (*data.PackageInfo, error) {
	pkg, err := app.RetrievePackage(pkgId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, nil
	}
	var pkgInfo = data.PackageInfo{}
	//pkgInfo.AddInfo(pkg, storehubPkg, true)
	enabled, err := app.GetServiceStatus(pkgId)
	if err != nil {
		return nil, errors.New("failed to check if package is enable " + pkgId)
	}
	pkgInfo.Enabled = enabled
	if task := installer.GetTask(pkgInfo.Id); task != nil {
		pkgInfo.Status = task.Status
	}
	return &pkgInfo, nil
}

// use to download and install app
func DownloadInstallApp(pkgId string, releaseType data2.ReleaseType) error {
	task, err := installer.CreateInstallTask(pkgId, releaseType)
	if err != nil {
		return err
	}
	task.Start()
	return nil
}

func UninstallApp(pkgId string) error {
	task := installer.CreateUninstallTask(pkgId)
	return task.Uninstall()
}

func CancelInstall(pkgId string) {
	installer.Cancel(pkgId)
}
