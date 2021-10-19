package protocol

import "github.com/cansulting/elabox-system-tools/foundation/event/data"

// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// this file provides funnction for generating response data for RPC
// response data is json string which will be encoded to Base64 before transmitting

// connectorClient.go
// interface for client communication and request handling
type ConnectorClient interface {
	GetState() data.ConnectionType
	// use to connect to local app server
	// @timeout: time in seconds it will timeout. @timeout > 0 to apply timeout
	Open(int16) error
	Close()
	// use to send request to event system
	SendSystemRequest(action string, data data.Action) (string, error)
	// use to subscribe to specific action.
	// @callback: will be called when someone broadcasted this action
	Subscribe(action string, callback interface{}) error
	// brodcast or send event to server
	Broadcast(event string, data interface{}) error
}
