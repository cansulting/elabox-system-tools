package main

import (
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/servicecenter"
	"log"
	"os"
	"time"
)

func main() {
	commandline := false
	if len(os.Args) > 1 {
		commandline = true
	}

	//global.Initialize()
	servicecenter.Initialize(commandline)
	defer servicecenter.Close()
	if err := appman.Initialize(commandline); err != nil {
		log.Panicln("installer failed to initialize " + err.Error())
		return
	}

	// this runs the server
	for global.Running {
		time.Sleep(time.Second * 1)
	}
	log.Println("System termination")
}
