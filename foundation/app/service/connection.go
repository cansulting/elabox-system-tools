package service

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
)

type ServiceDelegate func(message string, data interface{})

/*
	request.go
	2 way communication bridge between app and specific service.
*/
type Connection struct {
	ServiceId         string          // which service will be connected
	onServiceResponse ServiceDelegate // callback when recieved responses
	connector         protocol.ConnectorClient
}

// constructor for connection.
func NewConnection(
	connector protocol.ConnectorClient,
	serviceId string,
	onServiceResponse ServiceDelegate) (*Connection, error) {
	// listens to service responses
	if _, err := connector.Subscribe(
		data.Subscription{Action: serviceId}, onServiceResponse); err != nil {
		return nil, err
	}
	con := Connection{ServiceId: serviceId, onServiceResponse: onServiceResponse, connector: connector}
	return &con, nil
}

// sends specific request with data attached
func (t *Connection) RequestFor(action data.Action) (string, error) {
	return t.connector.SendServiceRequest(
		constants.SYSTEM_SERVICE_ID,
		action)
}
