// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package service

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
)

// callback function whenever recieve an action from server
type ServiceDelegate func(client protocol.ClientInterface, data data.Action) string

// 2 way communication bridge between app and specific service.
type RPCHandler struct {
	connector protocol.ConnectorClient
}

// constructor for RPCHandler.
func NewRPCHandler(connector protocol.ConnectorClient) *RPCHandler {
	con := RPCHandler{connector: connector}
	return &con
}

// use to listen to specific action from server, service delegate will be called upon response
func (t *RPCHandler) OnRecieved(action string, onServiceResponse ServiceDelegate) {
	// TODOserviceCommand := t.PackageId + ".service." + action
	t.connector.Subscribe(action, onServiceResponse)
}

// sends request to specific package with data attached
func (t *RPCHandler) Call(packageId string, action data.Action) (*data.Response, error) {
	strResponse, err := t.connector.SendServiceRequest(packageId, action)
	if err != nil {
		return nil, err
	}
	return &data.Response{Value: strResponse}, err
}

// send a request to system with data
func (t *RPCHandler) CallSystem(action data.Action) (*data.Response, error) {
	return t.Call(constants.SYSTEM_SERVICE_ID, action)
}

// use to broadcast to the system with specific action data
func (t *RPCHandler) CallBroadcast(action data.Action) (*data.Response, error) {
	return t.CallSystem(data.NewAction(constants.SYSTEM_BROADCAST, "", action))
}

// closes and uninitialize this handler
func (t *RPCHandler) Close() error {
	// t.connector.Broadcast(constants.SERVICE_UNBIND, nil)
	return nil
}
