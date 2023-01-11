package main

import (
	"encoding/base64"

	"github.com/cansulting/elabox-system-tools/foundation/account"
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
)

var pubkey []byte

type MyService struct {
}

func (instance *MyService) IsRunning() bool {
	return true
}

func (instance *MyService) OnStart() error {
	Controller.RPC.OnRecieved(AC_AUTH_DID, instance.onAuthDidAction)
	Controller.RPC.OnRecieved(AC_AUTH_SYSTEM, instance.onAuthSystem)
	Controller.RPC.OnRecieved(AC_SETUP_CHECK, instance.onCheckSetup)
	Controller.RPC.OnRecieved(AC_SETUP_DID, instance.onSetupDid)
	Controller.RPC.OnRecieved(constants.ACTION_RETRIEVE_PUBKEY, instance.onPublicKey)
	Controller.RPC.OnRecieved(AC_VALIDATE_TOKEN, instance.onValidateToken)
	return nil
}

func (instance *MyService) OnEnd() error {
	return nil
}

// authorize using system credentials
func (instance *MyService) onAuthSystem(client protocol.ClientInterface, action data.Action) string {
	// validate accounts
	params, err := action.DataToMap()
	if err != nil {
		rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed to parse parameters")
	}
	if params["username"] == nil || params["pass"] == nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, "either username or password is empty")
	}
	// authenticate account
	username := params["username"].(string)
	pass := params["pass"].(string)
	err, isValid := AuthenticateSystemAccount(username, pass)
	if err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	if !isValid {
		return rpc.CreateResponse(CREDENTIALS_INVALID, "username or password is invalid")
	}
	// create token
	token, err := CreateToken(username, "origin")
	if err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed generating token, "+err.Error())
	}
	ac := Account{
		Token:       token,
		DisplayName: "elabox",
		Status:      "active",
	}
	return rpc.CreateSuccessResponse(ac.ToJson())
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
	if _, valid := AuthenticateDid(did); !valid {
		return rpc.CreateResponse(rpc.INVALID_CODE, "incorrect did credentials")
	}
	acc := GetAccount()
	return rpc.CreateJsonResponse(rpc.SUCCESS_CODE, acc)
}

// use to check whether the elabox was setup already. return "setup" if already setup
func (instance *MyService) onCheckSetup(client protocol.ClientInterface, action data.Action) string {
	acc := GetAccount()
	if acc != nil && acc.Did != "" {
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
	acc := GetAccount()
	if acc.Did != "" {
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
			return rpc.CreateResponse(rpc.INVALID_PARAMETER_PROVIDED, "password is incorrect")
		}
	}
	presentation := acData["presentation"].(map[string]interface{})
	// step: validate presentation
	if presentation["holder"] == nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "no holder provider in presentation")
	}
	did := presentation["holder"].(string)
	if err := UpdateDid(did); err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed to setup did, "+err.Error())
	}
	// step: update wallet address
	if presentation["esc"] == nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, "no esc wallet address provided")
	}
	if err := UpdateWalletAddress("esc", presentation["esc"].(string)); err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, "failed updating esc wallet, "+err.Error())
	}
	SaveAccount()
	return rpc.CreateSuccessResponse("success")
}

// called when rpc was requested for public key
func (instance *MyService) onPublicKey(client protocol.ClientInterface, action data.Action) string {
	if pubkey == nil {
		key, err := GetPublicKey()
		key = []byte(base64.StdEncoding.EncodeToString(key))
		pubkey = key
		if err != nil {
			rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
		}
	}
	return rpc.CreateSuccessResponse(string(pubkey))
}

// called whe rpc request to validate token
func (instance *MyService) onValidateToken(client protocol.ClientInterface, action data.Action) string {
	isValid, err := account.ValidateTokenFromAction(action)
	if err != nil {
		rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	validStr := "valid"
	if !isValid {
		validStr = "invalid"
	}
	return rpc.CreateSuccessResponse(validStr)
}
