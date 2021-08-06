package server

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
)

// use to handle system services such as subscription and broadvcasting
// @handler: will be called for further system service handling
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
func Subscribe(server protocol.ConnectorServer, client protocol.ClientInterface, service string) string {
	if service == "" {
		service = constants.SYSTEM_SERVICE_ID
	}
	if err := server.SubscribeClient(client, service); err != nil {
		return err.Error()
	}
	return "subscribed to " + service
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
	broadcaster := constants.SYSTEM_SERVICE_ID
	if action.PackageId == "" {
		broadcaster = action.PackageId
	}
	err := server.Broadcast(broadcaster, action.Id, action)
	if err != nil {
		return err.Error()
	} else {
		return "success"
	}
}
