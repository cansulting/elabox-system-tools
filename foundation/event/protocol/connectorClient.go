package protocol

import "ela/foundation/event/data"

/*
	connectorClient.go

*/
type ConnectorClient interface {
	GetState() data.ConnectionType
	Open() error
	Close()
	// use to send request to specific service
	SendServiceRequest(serviceId string, data data.Action) (string, error)
	// use to subscribe to specific action.
	// @callback: will be called when someone broadcasted this action
	Subscribe(subs data.Subscription, callback interface{}) (string, error)
	Broadcast(action data.Action) error
}
