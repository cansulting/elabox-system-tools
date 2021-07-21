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
	println("For commands type help")
	// commandline is true if this app will do nothing aside from
	// commandline requests
	commandline := false
	if len(os.Args) > 1 {
		commandline = true
		processCmdline()
		return
	}

	// step: skip if system already running
	if connectToSystem() != nil {
		log.Println("System already running.")
		return
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
	log.Println("System terminated")
}
