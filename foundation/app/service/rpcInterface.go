package service

import (
	"ela/foundation/event/data"
)

type RPCInterface interface {
	// use to broadcast to the system
	CallSystem(action data.Action) (*data.Response, error)
	// use to broadcast to specific package
	Call(packageId string, action data.Action) (*data.Response, error)
	Disconnect() error
	OnRecieved(event string, onServiceResponse ServiceDelegate)
}
