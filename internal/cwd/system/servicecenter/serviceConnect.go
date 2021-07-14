package servicecenter

import (
	"ela/foundation/event/protocol"
)

/*
	This struct connects the bridge between the service provider and consumer
*/
type ServiceConnect struct {
	PackageId string
	Provider  protocol.ClientInterface
	Connector protocol.ConnectorServer
}

// creates new instance of service connect
// @client:
func NewServiceConnect(
	packageId string,
	provider protocol.ClientInterface,
	connector protocol.ConnectorServer) *ServiceConnect {
	newConnect := &ServiceConnect{
		Provider:  provider,
		PackageId: packageId,
		Connector: connector,
	}
	connector.Subscribe(packageId, newConnect.onServe)
	return newConnect
}

func (c *ServiceConnect) onServe(consumer protocol.ClientInterface, method string, data interface{}) string {
	response, err := c.Connector.BroadcastTo(c.Provider, method, data)
	if err != nil {
		return ""
	}
	return response
}
