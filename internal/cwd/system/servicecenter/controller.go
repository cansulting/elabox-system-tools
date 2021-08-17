package servicecenter

import (
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/web"
	"ela/server"
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

// use to register a service
// @serviceId: the packageId or the service id
func RegisterService(serviceId string, callback interface{}) {
	global.Server.EventServer.Subscribe(serviceId, callback)
}
