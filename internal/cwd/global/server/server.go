package server

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
)

// use to handle system services
// @next: will be called for system service handling
func InitSystemService(server protocol.ConnectorServer, handler func(protocol.ClientInterface, data.Action) interface{}) {
	// handle client subscription
	server.Subscribe(constants.SYSTEM_SERVICE_ID, func(client protocol.ClientInterface, action data.Action) interface{} {
		switch action.Id {
		// client wants to broadcast an action
		case constants.SYSTEM_BROADCAST:
			dataAc, err := action.DataToActionData()
			if err != nil {
				return err.Error()
			}
			return Broadcast(server, dataAc)
		// client wants to subscribe to specific action
		case constants.ACTION_SUBSCRIBE:
			return Subscribe(server, client, action.DataToString())
		}
		if handler != nil {
			return handler(client, action)
		}
		return "unknown"
	})

}

// callback when a client want to subscribe to specific action
func Subscribe(server protocol.ConnectorServer, client protocol.ClientInterface, action string) string {
	if action == "" {
		action = constants.SYSTEM_SERVICE_ID
	}
	if err := server.SubscribeClient(client, action); err != nil {
		return err.Error()
	}
	return "subscribed to " + action
}

// use to broadcast to action
func Broadcast(server protocol.ConnectorServer, action data.Action) string {
	/*
		pks, err = RetrievePackagesWithBroadcast(action.Id)
		if err != nil {
			return err.Error()
		}
		for _, pk := range pks {
			launchPackage(action, pk)
		}*/
	err := server.Broadcast(constants.SYSTEM_SERVICE_ID, action.Id, action)
	if err != nil {
		return err.Error()
	} else {
		return "success"
	}
}
