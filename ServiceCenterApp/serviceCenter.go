package main

import (
	base "ela.services/Base"
	servicecenter "ela.services/ServiceCenterLib"
)

/*
	serviceCenter.go
	file that manages the service connection
*/

////////////////////////GLOBAL DEFINITIONS///////////////

var connector servicecenter.ConnectorServer

func RunServer() {
	connector = servicecenter.CreateServerConnector()
	connector.Open()

	// register this service
	registerService(base.ACTIONCENTER_ID, OnRecievedRequest)
	// client wants to subscribe to action
	connector.Subscribe(base.ACTION_SUBSCRIBE, onClientSubscribeAction)
	connector.Subscribe(base.ACTION_BROADCAST, onClientBroadcastAction)
}

/// this closes the server
func Close() {
	connector.Close()
}

////////////////////////PRIVATE DEFINITIONS///////////////

// callback when a client want to subscribe to specific action
func onClientSubscribeAction(client servicecenter.ClientInterface, action string) {
	connector.SubscribeClient(client, action)
}

// callback when a client wants to broadcast a specific action
func onClientBroadcastAction(
	client servicecenter.ClientInterface,
	action servicecenter.ActionData) {
	connector.Broadcast(action.Action, action.Action, action)
}

// use to register a service
// @serviceId: the packageId or the service id
func registerService(serviceId string, callback interface{}) {
	connector.Subscribe(serviceId, callback)
}
