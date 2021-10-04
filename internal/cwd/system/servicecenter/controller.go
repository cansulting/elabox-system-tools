package servicecenter

import (
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/web"
	"github.com/cansulting/elabox-system-tools/server"
)

/*
	serviceCenter.go
	file that manages the service connection
*/

////////////////////////GLOBAL DEFINITIONS///////////////

// initiate services
func Initialize(commandline bool) error {
	if commandline {
		return nil
	}
	global.Server = &server.Manager{}
	global.Server.OnSystemEvent = OnRecievedRequest
	global.Server.Setup()
	global.Server.ListenAndServe()
	webservice := &web.WebService{}
	if err := webservice.Start(); err != nil {
		return err
	}
	// start running all services
	return nil
}

/// this closes the server
func Close() {
	if global.Server != nil {
		global.Server.Stop()
	}
}
