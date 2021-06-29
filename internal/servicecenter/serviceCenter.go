package main

import (
	"ela/foundation/constants"
	"ela/foundation/event"
	event_data "ela/foundation/event/data"
	"ela/foundation/event/protocol"
)

/*
	serviceCenter.go
	file that manages the service connection
*/

////////////////////////GLOBAL DEFINITIONS///////////////

var connector protocol.ConnectorServer

func RunServer() {
	connector = event.CreateServerConnector()
	connector.Open()

	// register this service
	registerService(constants.SERVICE_CENTER_ID, OnRecievedRequest)
	// client wants to subscribe to action
	connector.Subscribe(constants.ACTION_SUBSCRIBE, onClientSubscribeAction)
	connector.Subscribe(constants.ACTION_BROADCAST, onClientBroadcastAction)
}

/// this closes the server
func Close() {
	connector.Close()
}

////////////////////////PRIVATE DEFINITIONS///////////////

// callback when a client want to subscribe to specific action
func onClientSubscribeAction(client protocol.ClientInterface, action string) {
	connector.SubscribeClient(client, action)
}

// callback when a client wants to broadcast a specific action
func onClientBroadcastAction(
	client protocol.ClientInterface,
	action event_data.Action) {
	connector.Broadcast(action.Action, action.Action, action)
}

// use to register a service
// @serviceId: the packageId or the service id
func registerService(serviceId string, callback interface{}) {
	connector.Subscribe(serviceId, callback)
}
