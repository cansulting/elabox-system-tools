// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package event

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
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
	//println("Subscribe " + service)
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
	broadcastTo := constants.SYSTEM_SERVICE_ID
	if action.PackageId != "" {
		broadcastTo = action.PackageId
	}
	err := s.Broadcast(broadcastTo, action.Id, action)
	if err != nil {
		return err.Error()
	} else {
		return "success"
	}
}
