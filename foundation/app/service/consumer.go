package service

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"log"
)

type ServiceDelegate func(message string, data interface{})

/*
	request.go
	2 way communication bridge between app and specific service.
*/
type Consumer struct {
	PackageId string // which service will be connected
	connector protocol.ConnectorClient
}

// constructor for Consumer.
func NewConsumer(
	connector protocol.ConnectorClient,
	packageId string) (*Consumer, error) {
	log.Println("service.Consumer: connecting to service", packageId)
	// listens to service responses
	// if err := connector.Subscribe(serviceId, onServiceResponse); err != nil {
	//	return nil, err
	//}
	// send request to server that it will connect to specific service
	if err := connector.Broadcast(constants.SERVICE_BIND, packageId); err != nil {
		return nil, err
	}
	con := Consumer{PackageId: packageId, connector: connector}
	return &con, nil
}

func (t *Consumer) On(event string, onServiceResponse ServiceDelegate) {
	serviceCommand := t.PackageId + ".service." + event
	t.connector.Subscribe(serviceCommand, onServiceResponse)
}

// sends specific request with data attached
func (t *Consumer) RequestFor(action data.Action) (*data.Response, error) {
	strResponse, err := t.connector.SendServiceRequest(
		constants.SYSTEM_SERVICE_ID,
		action)
	if err != nil {
		return nil, err
	}
	return &data.Response{Value: strResponse}, err
}

func (t *Consumer) Disconnect() error {
	// t.connector.Broadcast(constants.SERVICE_UNBIND, nil)
	return nil
}
