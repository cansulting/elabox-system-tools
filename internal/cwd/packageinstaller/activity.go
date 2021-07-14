package main

import (
	appd "ela/foundation/app/data"
	"ela/foundation/errors"
	"ela/foundation/event/data"
	"ela/foundation/path"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type activity struct {
	running bool
}

func (a *activity) IsRunning() bool {
	return a.running
}

func (a *activity) OnStart(action data.Action) error {
	a.running = true
	// step: validate action
	sourcePkg := action.DataToString()
	if sourcePkg == "" {
		a.finish("Unable to locate package " + sourcePkg)
		return nil
	}
	log.Println("Installing package @", sourcePkg)
	// step: load package info
	pkgi := appd.DefaultPackage()
	if err := pkgi.LoadFromPackage(sourcePkg); err != nil {
		a.finish("Unable to load package " + sourcePkg)
		return nil
	}
	// step: if requires restart then start commandline mode
	if pkgi.Restart {
		if err := a.startCommandlineMode(sourcePkg); err != nil {
			a.finish("Unable to start commandline mode " + sourcePkg + "\n " + err.Error())
			return nil
		}
		a.finish("")
		return nil
	}
	// step: start installing
	backup := pkgi.IsSystemPackage()
	install := installer{BackupEnabled: backup, SilentInstall: false}
	if err := install.Decompress(sourcePkg); err != nil {
		a.finish("Unable to install file " + err.Error())
		return nil
	}
	// step: register package
	if err := install.RegisterPackage(); err != nil {
		a.finish("Unable to register package " + err.Error())
		return nil
	}
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

// install via commandline mode
func (a *activity) startCommandlineMode(pkgPath string) error {
	log.Println("Commandline mode initiated")
	// step: create a copy of this application and start exec
	binPath := path.GetAppMain(
		appController.Config.PackageId,
		!appController.Config.IsSystemPackage())
	dest := path.GetCacheDir() + "/" + filepath.Base(binPath)
	log.Println("Cloning binary @ " + dest)
	bytes, err := ioutil.ReadFile(binPath)
	if err != nil {
		return errors.SystemNew("installer.startCommandlineMode Failed to copy installer binary.", err)
	}
	if err := ioutil.WriteFile(dest, bytes, 0770); err != nil {
		return errors.SystemNew("installer.startCommandlineMode Failed to copy installer binary.", err)
	}
	// step: run the copied binary
	cmd := exec.Command(dest, pkgPath)
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	println("ERROR " + err.Error())
	//}
	//println(string(out))

	if err := cmd.Start(); err != nil {
		return errors.SystemNew("installer.startCommandlineMode Failed to execute commandline binary.", err)
	}
	log.Println("Running commandline installer")
	time.Sleep(2 * time.Second)
	os.Exit(0)
	return nil
}
