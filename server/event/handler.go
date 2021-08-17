package event

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
)

// use to handle system services such as subscription and broadvcasting
// @handler: will be called for further system service handling
func (s *SocketIOServer) HandleSystemService(handler func(protocol.ClientInterface, data.Action) interface{}) {
	// handle client subscription
	s.Subscribe(constants.SYSTEM_SERVICE_ID, func(client protocol.ClientInterface, action data.Action) interface{} {
		switch action.Id {
		// client wants to broadcast an action
		case constants.SYSTEM_BROADCAST:
			dataAc, err := action.DataToActionData()
			if err != nil {
				return err.Error()
			}
			return s.BroadcastAction(dataAc)
		// client wants to subscribe to specific action
		case constants.ACTION_SUBSCRIBE:
			return s.SubscribeToService(client, action.DataToString())
		}
		if handler != nil {
			return handler(client, action)
		}
		return "unknown"
	})

}

// callback when a client want to subscribe to specific action
func (s *SocketIOServer) SubscribeToService(client protocol.ClientInterface, service string) string {
	if service == "" {
		service = constants.SYSTEM_SERVICE_ID
	}
	if err := s.SubscribeClient(client, service); err != nil {
		return err.Error()
	}
	return "subscribed to " + service
}

// use to broadcast to action
func (s *SocketIOServer) BroadcastAction(action data.Action) string {
	/*
		pks, err = RetrievePackagesWithBroadcast(action.Id)
		if err != nil {
			return err.Error()
		}
		for _, pk := range pks {
			launchPackage(action, pk)
		}*/
	broadcaster := constants.SYSTEM_SERVICE_ID
	if action.PackageId != "" {
		broadcaster = action.PackageId
	}
	err := s.Broadcast(broadcaster, action.Id, action)
	if err != nil {
		return err.Error()
	} else {
		return "success"
	}
}
