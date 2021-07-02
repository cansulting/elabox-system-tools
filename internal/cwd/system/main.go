package main

import (
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/servicecenter"
	"log"
	"os"
	"time"
)

var running bool = true

func main() {
	commandline := false
	if len(os.Args) > 1 {
		commandline = true
	}

	global.Initialize()
	servicecenter.Initialize(commandline)
	defer servicecenter.Close()
	appman.Initialize(commandline)

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
					log.Panicln(err)
					return
				}
				log.Println("Success, ", data.PackageId, "was registered!")
				return
			}
			log.Panicln("Requires path to package installation dir")
		default:
			log.Panicln("Unsupported command " + command)
		}
		return
	}

	// this runs the server
	for running {
		time.Sleep(time.Second * 1)
	}
}
