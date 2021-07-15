package main

import (
	"C"
	"ela/foundation/app"
	"ela/internal/cwd/packageinstaller/global"
	"ela/internal/cwd/packageinstaller/logging"
	"log"
	"os"
)

func main() {
	InitializePath()
	logging.Initialize("installer.log")
	// install via commandline
	args := os.Args
	if len(args) > 1 {
		startCommandline()
		return
	}
	// install via installer service
	var err error
	global.AppController, err = app.NewController(&activity{}, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if err := app.RunApp(global.AppController); err != nil {
		log.Fatal(err)
	}
}
