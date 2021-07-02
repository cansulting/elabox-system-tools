package main

import (
	"ela/foundation/app"
	"log"
	"os"
)

var appController *app.Controller

func main() {
	// install via commandline
	args := os.Args
	if len(args) > 1 {
		// check if valid path
		packagePath := args[1]
		if _, err := os.Stat(packagePath); err != nil {
			log.Fatal("Unable to install package with invalid path.")
			return
		}
		// start install
		newInstall := installer{BackupEnabled: true, SilentInstall: false}
		if err := newInstall.decompress(packagePath); err != nil {
			log.Fatal(err.Error())
			return
		}
		// register package
		if err := newInstall.registerPackage(); err != nil {
			log.Fatal(err.Error())
			return
		}
		return
	}
	// install via installer service
	var err error
	appController, err = app.NewController(nil, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	app.RunApp(appController)
}
