package servicecenter

import (
	"os"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/appman"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"

	"log"

	"github.com/cansulting/elabox-system-tools/registry/app"
)

// callback when recieved service requests from client
func OnRecievedRequest(
	client protocol.ClientInterface,
	action data.Action,
) interface{} {
	log.Println("app.onRecievedRequest action=", action.Id)
	switch action.Id {
	case constants.ACTION_RPC:
		valAction, err := action.DataToActionData()
		if err != nil {
			return err.Error()
		}
		return sendPackageRPC(valAction)

	case constants.ACTION_START_ACTIVITY:
		activityAc, err := action.DataToActionData()
		if err != nil {
			return err.Error()
		}
		res := startActivity(activityAc, client)
		log.Println(res)
		return res
	case constants.APP_CHANGE_STATE:
		return onAppChangeState(client, action)
	case constants.SYSTEM_ACTIVITY_RESULT:
		return onReturnActivityResult(action)
	case constants.SYSTEM_UPDATE_MODE:
		return activateUpdateMode(client, action)
	case constants.SYSTEM_TERMINATE:
		return terminate(5)
	case constants.SYSTEM_TERMINATE_NOW:
		return terminate(0)
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

// client request to launch an activity
func startActivity(action data.Action, client protocol.ClientInterface) string {
	if action.Id == "" {
		return "appman.startActivity failed. No action was provided"
	}
	packageId := action.PackageId
	// step: look up for package id
	if packageId == "" {
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
	log.Println("Start activity", action.Id, packageId, action.DataToString())
	return `{"code":200, "message":"Launched"}`
}

// when activity returns a result
func onReturnActivityResult(action data.Action) string {
	packageId := action.PackageId
	if packageId == "" {
		return `{"code":400, "message": "package id should be the activity whho return the result"}`
	}
	app := appman.GetAppConnect(packageId, nil)
	if app == nil || app.StartedBy == "" {
		return `{"code":401, "message": "Cant find who app who will recieve result"}`
	}
	if !appman.IsAppRunning(app.StartedBy) {
		return `{"code":401, "message": "cant find the app that would recieve result"}`
	}
	originApp := appman.GetAppConnect(app.StartedBy, nil)
	return originApp.RPC.CallAct(action)
}

// send RPC to specific package
func sendPackageRPC(action data.Action) string {
	app := appman.GetAppConnect(action.PackageId, nil)
	if app == nil {
		return `{"code":401, "message": "Cant find package"}`
	}
	return app.RPC.CallAct(action)
}

// client requested to activate update mode
func activateUpdateMode(client protocol.ClientInterface, action data.Action) string {
	pk := action.DataToString()
	global.Server.EventServer.BroadcastAction(data.NewActionById(constants.BCAST_TERMINATE_N_UPDATE))
	startActivity(data.NewAction(constants.ACTION_APP_SYSTEM_INSTALL, "", pk), nil)
	return "success"
}

func terminate(seconds uint) string {
	go func() {
		log.Println("System will terminate after", seconds, "seconds")
		if seconds > 0 {
			time.Sleep(time.Second * time.Duration(seconds))
		}
		appman.TerminateAllApp()
		global.Server.EventServer.SetStatus(system.STOPPED, nil)
		global.Running = false
		time.Sleep(time.Millisecond * 100)
		os.Exit(0)
	}()
	return "success"
}
