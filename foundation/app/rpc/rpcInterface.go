// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package rpc

import (
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
)

type RPCInterface interface {
	// use to broadcast to the system
	CallSystem(action data.Action) (*Response, error)
	// use to broadcast to specific package
	CallRPC(packageId string, action data.Action) (*Response, error)
	StartActivity(action data.Action) (*Response, error)
	Close() error
	// set callback when recieved specific action/event
	OnRecieved(action string, onServiceResponse ServiceDelegate)
}
