package appman

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/internal/cwd/system/global"
)

// callback when recieved service requests from client
func OnRecievedRequest(
	client protocol.ClientInterface,
	action data.Action,
) interface{} {
	switch action.Id {
	// client wants to broadcast an action
	case constants.ACTION_BROADCAST:
		return onActionResolver(action.DataToActionData())
	// client wants to subscribe to specific action
	case constants.ACTION_SUBSCRIBE:
		return onClientSubscribeAction(client, action.DataToString())
	case constants.APP_CHANGE_STATE:
		return onAppChangeState(client, action)
	}
	return ""
}

func onAppChangeState(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	state := action.DataToActionData().DataToInt()

	// if awake then send pending actions to client
	if state == constants.APP_AWAKE {
		app := getAppConnect(action.PackageId)
		return app.pendingActions
	} else if state == constants.APP_SLEEP {
		// if sleep then wait to terminate the app
		removeAppConnect(action.PackageId, false)
	}
	return ""
}

// callback when a client want to subscribe to specific action
func onClientSubscribeAction(client protocol.ClientInterface, action string) string {
	if err := global.Connector.SubscribeClient(client, action); err != nil {
		return err.Error()
	}
	return ""
}
