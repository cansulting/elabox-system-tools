package service

import (
	"ela/foundation/event/data"
)

type RPCInterface interface {
	CallSystem(action data.Action) (*data.Response, error)
	Call(packageId string, action data.Action) (*data.Response, error)
	Disconnect() error
	OnRecieved(event string, onServiceResponse ServiceDelegate)
}
