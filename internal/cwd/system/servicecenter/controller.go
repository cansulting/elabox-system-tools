package servicecenter

import (
	"ela/foundation/constants"
	"ela/foundation/event"
	event_data "ela/foundation/event/data"
	"ela/foundation/event/protocol"
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

	// register this service
	RegisterService(constants.SYSTEM_SERVICE_ID, OnRecievedRequest)
	// client wants to subscribe to action
	//connector.Subscribe(constants.ACTION_SUBSCRIBE, onClientSubscribeAction)
	//connector.Subscribe(constants.ACTION_BROADCAST, onClientBroadcastAction)
}

/// this closes the server
func Close() {
	global.Connector.Close()
}

// use to register a service
// @serviceId: the packageId or the service id
func RegisterService(serviceId string, callback interface{}) {
	global.Connector.Subscribe(serviceId, callback)
}

////////////////////////PRIVATE DEFINITIONS///////////////
// callback when a client wants to broadcast a specific action
func onClientBroadcastAction(
	client protocol.ClientInterface,
	action event_data.Action) {
	global.Connector.Broadcast(action.Id, action.Id, action)
}
