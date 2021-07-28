package appman

import (
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"log"
)

/*
   This struct connects the bridge between the service client and consumer
*/
type RPCBridge struct {
	PackageId string                   // the target package
	Client    protocol.ClientInterface // the client of package
	Connector protocol.ConnectorServer
}

// creates new instance of service connect
// @client:
func NewRPCBridge(
	packageId string,
	client protocol.ClientInterface,
	connector protocol.ConnectorServer) *RPCBridge {
	newConnect := &RPCBridge{
		Client:    client,
		PackageId: packageId,
		Connector: connector,
	}
	connector.Subscribe(packageId, newConnect.onBridge)
	return newConnect
}

func (c *RPCBridge) onBridge(consumer protocol.ClientInterface, data data.Action) string {
	return c.CallAct(data)
}

// communicate to current package
func (c *RPCBridge) CallAct(data data.Action) string {
	if c.Client == nil {
		return "Ignored no client connected."
	}
	response, err := c.Connector.BroadcastTo(c.Client, data.Id, data)
	if err != nil {
		log.Println("RPCBridge.Call Response from", c.PackageId, "error "+err.Error())
		return err.Error()
	}
	return response
}

// call the owning package
func (c *RPCBridge) Call(action string, _data interface{}) string {
	return c.CallAct(data.NewAction(action, "", _data))
}
