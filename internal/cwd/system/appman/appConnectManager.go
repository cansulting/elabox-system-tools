// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file handles all currently running app via App Connect
// Use this class to run, stop and get running status of app

// This file manages the app connect.

package appman

import (
	"errors"

	appd "github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
	registry "github.com/cansulting/elabox-system-tools/registry/app"
)

// currently running processes
var running map[string]*AppConnect = make(map[string]*AppConnect)

func GetAllRunningApps() map[string]*AppConnect {
	return running
}

// use to run process for specific package. return true if success, false if already running
// client: the app's client
func GetAppConnect(packageId string, client protocol.ClientInterface) *AppConnect {
	app, ok := running[packageId]

	// is already running? return false
	if ok {
		if client != nil {
			app.Client = client
		}
		return app
	}
	// retrieve if already exist
	pk, _ := registry.RetrievePackage(packageId)
	if pk == nil {
		return nil
	}
	return AddAppConnect(pk, client)
}

// add app connect to list of running apps
func AddAppConnect(pk *appd.PackageConfig, client protocol.ClientInterface) *AppConnect {
	// create service if exist
	//var service *ServiceConnect = nil
	//if pk.HasServices() {
	//service = onServiceOpen(client, pk.PackageId)
	//}
	app := newAppConnect(pk, client)
	running[pk.PackageId] = app
	return app
}

// get package from running list
func LookupAppConnect(packageId string) *AppConnect {
	pk, ok := running[packageId]
	// is already running? return false
	if ok {
		return pk
	}
	return nil
}

// use to check if app is currently running or not
func IsAppRunning(packageId string) bool {
	app := LookupAppConnect(packageId)
	return app != nil
}

func RemoveAppConnect(packageId string, terminate bool) {
	app := LookupAppConnect(packageId)
	if app != nil {
		if terminate {
			if err := app.Terminate(); err != nil {
				global.Logger.Error().Err(err).Stack().Msg("Failed terminate " + app.PackageId + ". Trying force terminate.")
				if err := app.ForceTerminate(); err != nil {
					global.Logger.Error().Err(err).Caller().Msg("appConnectManager.TerminateAllApp failed force terminate ")
				}
			}
		}
		// close service
		// if app.Service != nil {
		// 	onServiceClose(packageId)
		// }
		delete(running, packageId)
	}
}

func TerminateAllApp() {
	global.Logger.Info().Msg("appConnectManager.TerminateAllApp() started")
	running := GetAllRunningApps()
	for pkid := range running {
		RemoveAppConnect(pkid, true)
	}
}

// this launches the activity
func LaunchAppActivity(
	packageId string,
	caller protocol.ClientInterface,
	pendingActivity data.Action) error {
	appc := GetAppConnect(packageId, nil)
	if appc == nil {
		return errors.New("package " + packageId + " is not installed")
	}
	if !appc.Config.HasActivity(pendingActivity.Id) {
		return errors.New("package " + packageId + " doesnt have a registered activity")
	}
	_, err := SendAppPendingAction(appc, pendingActivity, data.Action{})
	return err
}

func LaunchAppService(pkgid string) (*AppConnect, error) {
	app := GetAppConnect(pkgid, nil)
	if app == nil {
		return nil, errors.New("Package " + pkgid + " was not found.")
	}
	if app.launched {
		return app, nil
	}
	ac := data.NewActionById(constants.ACTION_START_SERVICE)
	app.PendingActions.AddPendingService(&ac)
	return app, app.Launch()
}

// use to launch app
func SendAppPendingAction(
	app *AppConnect,
	activityPending data.Action,
	servicePending data.Action) (*AppConnect, error) {

	if activityPending.Id != "" {
		app.PendingActions.AddPendingActivity(&activityPending)
	}
	if servicePending.Id != "" {
		app.PendingActions.AddPendingService(&servicePending)
	}

	return app, app.Launch()
}

// run all start up apps
func InitializeStartups() {
	global.Logger.Info().Msg("Services are starting up...")
	pkgs, err := registry.RetrieveStartupPackages()
	if err != nil {
		global.Logger.Error().Err(err).Caller().Msg("Failed retrieving startup packages.")
	}
	for _, pkg := range pkgs {
		if _, err := LaunchAppService(pkg.PackageId); err != nil {
			global.Logger.Error().Err(err).Caller().Msg("Failed launching app.")
		}
	}
}
