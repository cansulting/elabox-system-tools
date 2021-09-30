package main

import (
	"ela/foundation/system"
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/config"
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/servicecenter"
	"os"
	"time"
)

func main() {
	// commandline is true if this app will do nothing aside from
	// commandline requests
	commandline := false
	if len(os.Args) > 1 {
		commandline = true
		processCmdline()
		return
	}
	println("For commands type help")

	// step: skip if system already running
	if connectToSystem() != nil {
		println("System already running.")
		return
	}
	if err := config.Init(); err != nil {
		global.Logger.Panic().Err(err).Caller().Msg("Failed initializing config.")
		return
	}
	global.Logger.Info().Msg("System start running...")
	servicecenter.Initialize(commandline)
	global.Server.EventServer.SetStatus(system.BOOTING, nil)
	defer servicecenter.Close()
	if err := appman.Initialize(commandline); err != nil {
		global.Logger.Panic().Err(err).Caller().Msg("Application manager failed to initialize.")
		return
	}
	global.Server.EventServer.SetStatus(system.RUNNING, nil)
	// this runs the server
	for global.Running {
		time.Sleep(time.Second * 1)
	}
	global.Logger.Info().Msg("System terminated")
}
