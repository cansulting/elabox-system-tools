package service

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"log"
)

/*
	request.go
	2 way communication bridge between app and specific service.
*/
type Provider struct {
	PackageId string // which service will be connected
	connector protocol.ConnectorClient
}

// constructor for Provider.
func NewProvider(
	connector protocol.ConnectorClient,
	packageId string) (*Provider, error) {
	log.Println("service.Provider: initiating service", packageId)
	// send request to server that it will connect to specific service
	if err := connector.Broadcast(constants.SERVICE_BIND, packageId); err != nil {
		return nil, err
	}
	con := Provider{PackageId: packageId, connector: connector}
	return &con, nil
}

// someone request for service
func (t *Provider) OnServe(event string, onServiceResponse ServiceDelegate) {
	serviceCommand := t.PackageId + ".service." + event
	t.connector.Subscribe(serviceCommand, onServiceResponse)
}

// send to all consumer
func (t *Provider) ServeToAll(action data.Action) (*data.Response, error) {
	strResponse, err := t.connector.SendServiceRequest(
		constants.SYSTEM_SERVICE_ID,
		action)
	if err != nil {
		return nil, err
	}
	return &data.Response{Value: strResponse}, err
}

func (t *Provider) Disconnect() error {
	return t.connector.Broadcast(constants.SERVICE_UNBIND, nil)
}
