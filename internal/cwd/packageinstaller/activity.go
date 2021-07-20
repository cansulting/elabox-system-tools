package main

import (
	appd "ela/foundation/app/data"
	"ela/foundation/event/data"
	"ela/internal/cwd/packageinstaller/global"
	"log"
)

type activity struct {
	running bool
}

func (a *activity) IsRunning() bool {
	return a.running
}

func (a *activity) OnStart(action data.Action) error {
	a.running = true
	global.AppController.SetActivityResult(0.0)
	// step: validate action
	sourcePkg := action.DataToString()
	if sourcePkg == "" {
		a.finish("Unable to locate package " + sourcePkg)
		return nil
	}
	log.Println("Installing package @", sourcePkg)
	// step: load package info
	pkgi := appd.DefaultPackage()
	if err := pkgi.LoadFromZipPackage(sourcePkg); err != nil {
		a.finish("Unable to load package " + sourcePkg)
		return nil
	}
	// step: start installing
	backup := pkgi.IsSystemPackage()
	install := installer{BackupEnabled: backup}
	if err := install.Decompress(sourcePkg); err != nil {
		a.finish("Unable to install file " + err.Error())
		return nil
	}
	// step: register package
	if err := install.Postinstall(); err != nil {
		a.finish("Unable to register package " + err.Error())
		return nil
	}
	global.AppController.SetActivityResult(1.0)
	a.finish("")
	return nil
}

func (a *activity) OnEnd() error {
	return nil
}

func (a *activity) finish(err string) {
	if err != "" {
		log.Println(err)
	} else {
		log.Println("Install success")
	}
	a.running = false
}
