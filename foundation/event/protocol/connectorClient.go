package protocol

import "ela/foundation/event/data"

/*
	connectorClient.go

*/
type ConnectorClient interface {
	GetState() data.ConnectionType
	// use to connect to local app server
	// @timeout: time in seconds it will timeout. @timeout > 0 to apply timeout
	Open(int16) error
	Close()
	// use to send request to specific service
	SendServiceRequest(serviceId string, data data.Action) (string, error)
	// use to subscribe to specific action.
	// @callback: will be called when someone broadcasted this action
	Subscribe(action string, callback interface{}) error
	Broadcast(event string, data interface{}) error
}
