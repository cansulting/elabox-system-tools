package servicecenter

import (
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	"ela/internal/cwd/system/appman"
	"ela/internal/cwd/system/global"
	"ela/internal/cwd/system/system_update"

	"ela/registry/app"
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
	case constants.SYSTEM_BROADCAST:
		return processBroadcastAction(action.DataToActionData())
	// client wants to subscribe to specific action
	case constants.ACTION_SUBSCRIBE:
		return onClientSubscribeAction(client, action.DataToString())
	case constants.APP_CHANGE_STATE:
		return onAppChangeState(client, action)
	case constants.SYSTEM_ACTIVITY_RESULT:
		return onReturnActivityResult(action)
	case constants.SYSTEM_UPDATE_MODE:
		return activateUpdateMode(client, action)
	case constants.SYSTEM_TERMINATE:
		return terminate()
	}
	return ""
}

func onAppChangeState(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	state := action.DataToInt()

	// if awake then send pending actions to client
	if state == constants.APP_AWAKE {
		app := appman.GetAppConnect(action.PackageId, client)
		if app != nil {
			return app.PendingActions
		} else {
			log.Println("appman.OnAppchangeState() package was not registered for ",
				action.PackageId,
				"App component might not work properly.")
			return nil
		}
	} else if state == constants.APP_SLEEP {
		// if sleep then wait to terminate the app
		appman.RemoveAppConnect(action.PackageId, false)
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

// client request to launch an activity
func startActivity(action data.Action, client protocol.ClientInterface) string {
	if action.Id == "" {
		return "appman.startActivity failed. No action was provided"
	}
	packageId := action.PackageId
	// step: look up for package id
	if packageId != "" {
		pks, err := app.RetrievePackagesWithActivity(action.Id)
		if err != nil {
			return errors.SystemNew("appman.StartActivity failed to start "+packageId, err).Error()
		}
		if len(pks) > 0 {
			packageId = pks[0]
		} else {
			return "cant find package with action " + action.Id
		}
	}
	if err := appman.LaunchAppActivity(packageId, client, action); err != nil {
		return `{"code":401, "message":` + err.Error() + ` }`
	}
	return `{"code":200, "message":"Launched"}`
}

// when activity returns a result
func onReturnActivityResult(action data.Action) string {
	packageId := action.PackageId
	if packageId == "" {
		return `{"code":400, "message": "package id should be the activity whho return the result"}`
	}
	if !appman.IsAppRunning(packageId) {
		return `{"code":401, "message": "cant find the app that would recieve result"}`
	}
	app := appman.GetAppConnect(packageId, nil)
	if app.StartedBy == "" {
		return `{"code":401, "message": "Cant find who app who will recieve result"}`
	}
	if !appman.IsAppRunning(app.StartedBy) {
		return `{"code":401, "message": "cant find the app that would recieve result"}`
	}
	originApp := appman.GetAppConnect(app.StartedBy, nil)
	return originApp.RPC.CallAct(action)
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

// client requested to activate update mode
func activateUpdateMode(client protocol.ClientInterface, action data.Action) string {
	pkconfig, err := app.RetrievePackage(global.INSTALLER_PKG_ID)
	if err != nil {
		log.Println("Package retrieve error", err.Error())
		return err.Error()
	}
	system_update.Start(action.DataToString(), pkconfig)
	global.Running = false
	return "success"
}

func terminate() string {
	appman.TerminateAllApp()
	global.Running = false
	return "success"
}