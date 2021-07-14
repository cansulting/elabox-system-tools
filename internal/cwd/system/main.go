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

	// process commandlines
	if commandline {
		command := os.Args[1]
		commandLen := len(command) - 1
		switch command {
		case "package-reg":
			if commandLen >= 2 {
				dir := os.Args[2]
				data, err := appman.RegisterPackageSrc(dir)
				if err != nil {
					log.Fatalln(err.Error())
					return
				}
				log.Println("Success, ", data.PackageId, "was registered!")
				return
			}
			log.Fatalln("Requires path to package installation dir")
		default:
			log.Fatalln("Unsupported command " + command)
		}
		return
	}

	// this runs the server
	for global.Running {
		time.Sleep(time.Second * 1)
	}
	log.Println("System termination")
}
