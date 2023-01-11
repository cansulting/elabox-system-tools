// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// protocol for action server

package protocol

import (
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/system"
)

// interface for service communication to clients
type ConnectorServer interface {
	GetState() data.ConnectionType
	Open() error
	SetStatus(status system.Status) error
	GetStatus() system.Status
	// send data to all room
	Broadcast(room string, event string, data interface{}) error
	// send service response to client
	// @recipient the recipient of action
	BroadcastTo(recipient ClientInterface, method string, data interface{}) (string, error)
	// server listen to room
	Subscribe(room string, callback interface{}) error
	/// make the client listen to room
	SubscribeClient(socketClient ClientInterface, room string) error
	Close() error
}
