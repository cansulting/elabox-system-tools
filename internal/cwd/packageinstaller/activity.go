package main

import (
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	global "github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/pkg"
)

type activity struct {
	running bool
}

func (a *activity) IsRunning() bool {
	return a.running
}

func (a *activity) OnStart(action *data.Action) error {
	a.running = true
	// step: validate action
	sourcePkg := action.DataToString()
	global.Logger.Info().Msg("Installing package @" + sourcePkg)
	pkgData, err := pkg.LoadFromSource(sourcePkg)
	if err != nil {
		a.finish(err.Error())
		return nil
	}
	if action.Id == constants.ACTION_APP_INSTALL {
		return a.startNormalInstall(pkgData)
	}
	return a.startSystemInstall(sourcePkg, pkgData)
}

func (a *activity) startNormalInstall(pkgd *pkg.Data) error {
	// step: start installing
	backup := pkgd.Config.IsSystemPackage()
	install := NewInstaller(pkgd, backup)
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
func (a *activity) startSystemInstall(pkgSource string, pkgd *pkg.Data) error {
	if pkgd.HasCustomInstaller() {
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
	return a.startNormalInstall(pkgd)
}

func (a *activity) OnEnd() error {
	return nil
}

func (a *activity) finish(err string) {
	if err != "" {
		global.Logger.Error().Caller().Stack().Msg(err)
	} else {
		global.Logger.Info().Msg("Install success")
	}
	a.running = false
}
