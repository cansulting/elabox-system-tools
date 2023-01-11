package main

import (
	"os"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/appman"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/env"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/servicecenter"
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
	// step: skip if system already running
	if connectToSystem() != nil {
		println("System already running.")
		return
	}
	if err := env.Init(); err != nil {
		global.Logger.Panic().Err(err).Caller().Msg("Failed initializing config.")
		return
	}
	global.Logger.Info().Msg("System start running...")
	if err := servicecenter.Initialize(commandline); err != nil {
		global.Logger.Panic().Err(err).Caller().Msg("Failed initializing service center. " + err.Error())
		return
	}
	global.Server.EventServer.SetStatus(system.BOOTING)
	defer servicecenter.Close()
	if err := appman.Initialize(commandline); err != nil {
		global.Logger.Panic().Err(err).Caller().Msg("Application manager failed to initialize.")
		return
	}
	global.Server.EventServer.SetStatus(system.RUNNING)
	// this runs the server
	for global.Running {
		time.Sleep(time.Second * 1)
	}
	global.Logger.Info().Msg("System terminated")
}
