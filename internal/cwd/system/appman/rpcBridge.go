// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file provides two way communication between two apps.

package appman

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
)

/*
   This struct connects the bridge between the service client and consumer
*/
type RPCBridge struct {
	PackageId string      // the target package
	App       *AppConnect // the client of package
	Connector protocol.ConnectorServer
}

// creates new instance of service connect
// @client:
func NewRPCBridge(
	packageId string,
	app *AppConnect,
	connector protocol.ConnectorServer) *RPCBridge {
	newConnect := &RPCBridge{
		App:       app,
		PackageId: packageId,
		Connector: connector,
	}
	if err := connector.Subscribe(packageId, newConnect.onBridge); err != nil {
		global.Logger.Error().Err(err).Caller().Msg("Failed subscribing to " + packageId)
	}
	return newConnect
}

func (c *RPCBridge) onBridge(consumer protocol.ClientInterface, data data.Action) string {
	res, err := c.CallAct(data)
	if err != nil {
		rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse(res)
}

// communicate to current package
func (c *RPCBridge) CallAct(data data.Action) (string, error) {
	if c.App == nil {
		return "", errors.SystemNew("Ignored no client connected for "+c.PackageId, nil)
	}
	response, err := c.Connector.BroadcastTo(c.App.Client, data.Id, data)
	if err != nil {
		return "", err
	}
	return response, nil
}

// call the owning package
func (c *RPCBridge) Call(action string, _data interface{}) (string, error) {
	return c.CallAct(data.NewAction(action, "", _data))
}
