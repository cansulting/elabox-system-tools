package service

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
)

type ServiceDelegate func(client protocol.ClientInterface, data data.Action) string

/*
	request.go
	2 way communication bridge between app and specific service.
*/
type RPCHandler struct {
	connector protocol.ConnectorClient
}

// constructor for RPCHandler.
func NewRPCHandler(connector protocol.ConnectorClient) *RPCHandler {
	con := RPCHandler{connector: connector}
	return &con
}

func (t *RPCHandler) OnRecieved(action string, onServiceResponse ServiceDelegate) {
	// TODOserviceCommand := t.PackageId + ".service." + action
	t.connector.Subscribe(action, onServiceResponse)
}

// sends specific request with data attached
func (t *RPCHandler) Call(packageId string, action data.Action) (*data.Response, error) {
	strResponse, err := t.connector.SendServiceRequest(packageId, action)
	if err != nil {
		return nil, err
	}
	return &data.Response{Value: strResponse}, err
}

func (t *RPCHandler) CallSystem(action data.Action) (*data.Response, error) {
	return t.Call(constants.SYSTEM_SERVICE_ID, action)
}

// use to broadcast to the system
func (t *RPCHandler) CallBroadcast(action data.Action) (*data.Response, error) {
	return t.CallSystem(data.NewAction(constants.SYSTEM_BROADCAST, "", action))
}

func (t *RPCHandler) Disconnect() error {
	// t.connector.Broadcast(constants.SERVICE_UNBIND, nil)
	return nil
}
