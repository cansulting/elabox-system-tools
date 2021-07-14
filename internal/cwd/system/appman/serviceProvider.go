package appman

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/internal/cwd/system/global"
	"log"
)

// callback when recieved service requests from client
func OnRecievedRequest(
	client protocol.ClientInterface,
	action data.Action,
) interface{} {
	log.Println("app.onRecievedRequest action=", action.Id)
	switch action.Id {
	case constants.ACTION_START_ACTIVITY:
		res := startActivity(action.DataToActionData(), client)
		log.Println(res)
		return res
	// client wants to broadcast an action
	case constants.ACTION_BROADCAST:
		return processBroadcastAction(action.DataToActionData())
	// client wants to subscribe to specific action
	case constants.ACTION_SUBSCRIBE:
		return onClientSubscribeAction(client, action.DataToString())
	case constants.APP_CHANGE_STATE:
		return onAppChangeState(client, action)
	case constants.SYSTEM_UPDATE_MODE:
		return activateUpdateMode(client, action)
	}
	return ""
}

func onAppChangeState(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	state := action.DataToInt()

	// if awake then send pending actions to client
	if state == constants.APP_AWAKE {
		app := getAppConnect(action.PackageId, client)
		if app != nil {
			return app.pendingActions
		} else {
			return nil
		}
	} else if state == constants.APP_SLEEP {
		// if sleep then wait to terminate the app
		RemoveAppConnect(action.PackageId, false)
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

// launch for activitiy
func startActivity(action data.Action, client protocol.ClientInterface) string {
	if action.Id == "" {
		return "appman.startActivity failed. No action was provided"
	}
	packageId := action.PackageId
	// step: look up for package id
	if packageId != "" {
		pks, err := RetrievePackagesWithActivity(action.Id)
		if err != nil {
			return errors.SystemNew("appman.StartActivity failed to start "+packageId, err).Error()
		}
		if len(pks) > 0 {
			packageId = pks[0]
		} else {
			return "cant find package with action " + action.Id
		}
	}
	// start launching the activity
	app := getAppConnect(packageId, client)
	app.pendingActions.AddPendingActivity(action)
	err := app.Launch()
	if err != nil {
		return err.Error()
	}
	return ""
}

func processBroadcastAction(action data.Action) string {
	/*
		pks, err = RetrievePackagesWithBroadcast(action.Id)
		if err != nil {
			return err.Error()
		}
		for _, pk := range pks {
			launchPackage(action, pk)
		}*/
	return ""
}

func activateUpdateMode(client protocol.ClientInterface, action data.Action) string {
	log.Println("appman.activateUpdateMode() started")
	TerminateAllApp()
	global.Running = false
	return "success"
}
