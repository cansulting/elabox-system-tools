package main

import (
	"ela/foundation/app"
	"log"
	"os"
)

var appController *app.Controller

func main() {
	InitializePath()
	// install via commandline
	args := os.Args
	if len(args) > 1 {
		log.Println("Elabox Installer Commandline")
		// step: check if valid path
		packagePath := args[1]
		if _, err := os.Stat(packagePath); err != nil {
			log.Fatal("Unable to install package with invalid path.")
			return
		}
		// step: terminate system. to make sure we dont have issues in installing system related files
		newInstall := installer{BackupEnabled: true, SilentInstall: true}
		if err := newInstall.TerminateSystem(); err != nil {
			log.Fatalln(err)
			return
		}
		// step: start install
		if err := newInstall.Decompress(packagePath); err != nil {
			log.Fatal(err.Error())
			return
		}
		// step: register package
		if err := newInstall.RegisterPackage(); err != nil {
			log.Fatal(err.Error())
			return
		}
		// step: restart system
		if err := newInstall.RestartSystem(); err != nil {
			log.Fatal(err.Error())
			return
		}
		return
	}
	// install via installer service
	var err error
	appController, err = app.NewController(&activity{}, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if err := app.RunApp(appController); err != nil {
		log.Fatal(err)
	}
}
