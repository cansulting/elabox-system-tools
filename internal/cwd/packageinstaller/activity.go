// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// controller.go
// Controller class handles the application lifecycle.
// To initialize call NewController, for debugging use NewControllerWithDebug
// please see the documentation for more info.

// this file handles activity component of installer

package main

import (
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/broadcast"
	global "github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/pkg"
)

type activity struct {
	running    bool
	currentPkg string // current package id being installed
}

func (a *activity) IsRunning() bool {
	return a.running
}

// callback when the activity is started
func (a *activity) OnStart() error {
	a.running = true
	return nil
}

// callback when recieved a pending action from system
func (a *activity) OnPendingAction(action *data.Action) error {
	// step: validate action
	sourcePkg := action.DataToString()
	a.currentPkg = sourcePkg
	broadcast.UpdateSystem(a.currentPkg, broadcast.INITIALIZING)
	global.Logger.Info().Msg("Installing package @" + sourcePkg)
	pkgData, err := pkg.LoadFromSource(sourcePkg)
	if err != nil {
		a.finish(err.Error())
		return nil
	}
	if action.Id == constants.ACTION_APP_INSTALL ||
		!pkgData.HasCustomInstaller() {
		return a.startNormalInstall(pkgData)
	}
	return a.runCustomInstaller(sourcePkg, pkgData)
}

func (a *activity) startNormalInstall(pkgd *pkg.Data) error {
	// step: start installing
	backup := pkgd.Config.IsSystemPackage()
	install := NewInstaller(pkgd, backup)
	install.SetProgressListener(a.onInstallProgress)
	install.SetErrorListener(a.onInstallError)
	broadcast.UpdateSystem(pkgd.Config.PackageId, broadcast.INPROGRESS)
	if err := install.Start(); err != nil {
		a.finish("Unable to install file " + err.Error())
		return nil
	}
	// step: register package
	if err := install.Finalize(); err != nil {
		a.finish("Unable to register package " + err.Error())
		return nil
	}
	a.finish("")
	return nil
}

// system install
func (a *activity) runCustomInstaller(pkgSource string, pkgd *pkg.Data) error {
	// start custom installer
	if err := pkgd.RunCustomInstaller(pkgSource, false, "-s", "-l", "-i"); err != nil {
		return errors.SystemNew("Failed installing system package.", err)
	}
	a.running = false
	time.Sleep(time.Millisecond * 200)
	// system terminate
	global.AppController.RPC.CallSystem(data.NewActionById(constants.SYSTEM_TERMINATE_NOW))
	return nil
}

func (a *activity) OnEnd() error {
	return nil
}

func (a *activity) finish(err string) {
	if err != "" {
		global.Logger.Error().Caller().Stack().Msg(err)
		broadcast.Error(a.currentPkg, 0, err)
	} else {
		global.Logger.Info().Msg("Install success")
		broadcast.UpdateSystem(a.currentPkg, broadcast.SUCCESS)
	}
	// comment this line. system will terminate this activity automatically
	//a.running = false
}

// callback from installer on progress
func (a *activity) onInstallProgress(progress int, pkg string) {
	broadcast.UpdateProgress(pkg, progress)
}

// callback from installer on error
func (a *activity) onInstallError(pkg string, code int, reason string, err error) {
	broadcast.Error(pkg, code, reason)
}
