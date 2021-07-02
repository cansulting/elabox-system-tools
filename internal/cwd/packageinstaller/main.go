package main

import (
	"ela/foundation/app"
	"fmt"
	"log"
)

var appController *app.Controller

func main() {
	//backup := Backup{}
	//error := backup.LoadAndApply("system.backup")
	//print(error.Error())
	//return
	newInstall := installer{backupEnabled: true}
	newInstall.initInstall(false)
	err := newInstall.decompress("sample.ela")
	if err != nil {
		fmt.Println(err.Error())
	}
	appController, err = app.NewController(nil, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	app.RunApp(appController)
}
