package servicecenter

import (
	"ela/foundation/constants"
	"ela/foundation/event"
	"ela/internal/cwd/system/global"
)

/*
	serviceCenter.go
	file that manages the service connection
*/

////////////////////////GLOBAL DEFINITIONS///////////////

// initiate services
func Initialize(commandline bool) {
	if commandline {
		return
	}
	global.Connector = event.CreateServerConnector()
	global.Connector.Open()
	global.Connector.Subscribe(constants.SYSTEM_SERVICE_ID, OnRecievedRequest)
	// start running all services
}

/// this closes the server
func Close() {
	if global.Connector != nil {
		global.Connector.Close()
	}
}

// use to register a service
// @serviceId: the packageId or the service id
func RegisterService(serviceId string, callback interface{}) {
	global.Connector.Subscribe(serviceId, callback)
}
