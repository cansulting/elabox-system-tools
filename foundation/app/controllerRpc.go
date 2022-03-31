// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// this file handles app related rpc callbacks and events0

package app

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
)

func (instance *Controller) initRPCRequests() {
	instance.RPC.OnRecieved(constants.APP_TERMINATE, instance.onTerminate)
	instance.RPC.OnRecieved(constants.SERVICE_PENDING_ACTIONS, instance.onPendingActions)
}

// callback from system. this app requested to be terminated
func (instance *Controller) onTerminate(client protocol.ClientInterface, data data.Action) string {
	instance.End()
	return ""
}

// callback from system. when recieved a pending action
func (instance *Controller) onPendingActions(client protocol.ClientInterface, action data.Action) string {
	pendingAc, err := action.DataToMap()
	if err != nil {
		logger.GetInstance().Error().Err(err).Caller().Msg("failed to parse pending actions")
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	actionG := data.NewActionGroupFromMap(pendingAc)
	// step: parse activity value
	if actionG.Activity != nil {
		// forward to activity
		if err := instance.Activity.OnPendingAction(actionG.Activity); err != nil {
			logger.GetInstance().Error().Err(err).Caller().Msg("failed processing pending action")
			return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
		}
	}

	return rpc.CreateSuccessResponse("success")
}
