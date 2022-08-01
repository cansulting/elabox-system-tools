package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
)

type MyService struct {
}

func (instance *MyService) IsRunning() bool {
	return true
}

func (instance *MyService) OnStart() error {
	Controller.RPC.OnRecieved(AC_DEVICE_SERIAL, instance.onGetDeviceSerial)
	Controller.RPC.OnRecieved(AC_AUTH_DID, instance.onAuthDidAction)
	Controller.RPC.OnRecieved(AC_SETUP_CHECK, instance.onCheckSetup)
	Controller.RPC.OnRecieved(AC_SETUP_DID, instance.onSetupDid)
	return nil
}

func (instance *MyService) OnEnd() error {
	return nil
}

// authorize user given the did presentation, upon success returns JWT token
func (instance *MyService) onAuthDidAction(client protocol.ClientInterface, action data.Action) string {
	presentation, err := action.DataToMap()
	// step: validate presentation
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, "invalid did presentation provided, "+err.Error())
	}
	if presentation["holder"] == nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, "no holder was provided")
	}
	// step: compare with the existing hash
	did := presentation["holder"].(string)
	if AuthenticateDid(did) {
		acc := RetrievePrimaryAccountDetails()
		return rpc.CreateJsonResponse(rpc.SUCCESS_CODE, acc)
	} else {
		return rpc.CreateResponse(rpc.INVALID_CODE, "incorrect did credentials")
	}
}

// use to check whether the elabox was setup already. return "setup" if already setup
func (instance *MyService) onCheckSetup(client protocol.ClientInterface, action data.Action) string {
	if IsDidSetup() {
		return rpc.CreateSuccessResponse("setup")
	} else {
		return rpc.CreateSuccessResponse("not_setup")
	}
}

// use to setup did, requires did presentation, username and password
func (instance *MyService) onSetupDid(client protocol.ClientInterface, action data.Action) string {
	acData, err := action.DataToMap()
	// step: validate inputs
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_PARAMETER_PROVIDED, "invalid parameters provided, "+err.Error())
	}
	if acData["presentation"] == nil {
		return rpc.CreateResponse(rpc.INVALID_PARAMETER_PROVIDED, "no presentation provided")
	}

	// step: authenticate for existing did setup
	if IsDidSetup() {
		if acData["username"] == nil || acData["password"] == nil {
			return rpc.CreateResponse(rpc.INVALID_PARAMETER_PROVIDED, "username and password is required")
		}
		pass := acData["password"].(string)
		username := acData["username"].(string)
		err, success := AuthenticateSystemAccount(username, pass)
		if err != nil {
			return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
		}
		if !success {
			return rpc.CreateResponse(rpc.INVALID_PARAMETER_PROVIDED, "password or username is incorrect")
		}
	}
	presentation := acData["presentation"].(map[string]interface{})
	if err := SetDeviceDid(presentation); err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed to setup did, "+err.Error())
	}
	return rpc.CreateSuccessResponse("success")
}
