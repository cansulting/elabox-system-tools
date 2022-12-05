package main

import (
	"errors"
	"sort"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/installer"
	"github.com/cansulting/elabox-system-tools/registry/app"
)

var deviceSerial = ""

// retrieve all apps
// @beta is true if include all apps for testing and demo apps
func RetrieveAllApps(beta bool) ([]data.PackageInfo, error) {
	// step: retrieve all apps from registry
	storeItems, err := app.RetrieveAllPackages()
	if err != nil {
		return nil, errors.New("unable to retrieve all installed packages. inner: " + err.Error())
	}
	var previews = make([]data.PackageInfo, 0, len(storeItems))
	var tmpPreview data.PackageInfo
	// step: iterate on packages
	for _, pkg := range storeItems {
		installedInfo, err := app.RetrievePackage(pkg)
		if err != nil {
			logger.GetInstance().Debug().Msg("unable to retrieve cache item for package: " + pkg + ". inner: " + err.Error())
		}
		tmpPreview = data.PackageInfo{}
		tmpPreview.AddInfo(installedInfo, nil, false)
		if task := installer.GetTask(tmpPreview.Id); task != nil {
			tmpPreview.Status = task.Status
		}
		// check if currently in download
		previews = append(previews, tmpPreview)
	}
	// add tasks installing
	for _, task := range installer.GetAllTasks() {
		if task.Status == global.Downloading || task.Status == global.Installing {
			tmp := task.Definition.ToPackageInfo()
			previews = append(previews, tmp)
		}
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
func DownloadInstallApp(link data.InstallDef, dependencies []data.InstallDef) error {
	task, err := installer.CreateInstallTask(link, dependencies)
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
