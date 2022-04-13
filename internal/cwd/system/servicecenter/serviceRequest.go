// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// This file provide services for client requests for RPC, broadcast, starting activity etc.

package servicecenter

import (
	"os"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/appman"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/config"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/debugging"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"

	"log"

	"github.com/cansulting/elabox-system-tools/registry/app"
)

// callback when recieved service requests from client
func OnRecievedRequest(
	client protocol.ClientInterface,
	action data.Action,
) interface{} {
	if config.GetBuildMode() == config.DEBUG {
		logger.GetInstance().Debug().Msg("onRecievedRequest action=" + action.Id)
	}
	switch action.Id {
	case constants.ACTION_RPC:
		valAction, err := action.DataToActionData()
		if err != nil {
			return err.Error()
		}
		return sendPackageRPC(action.PackageId, valAction)

	case constants.ACTION_START_ACTIVITY:
		activityAc, err := action.DataToActionData()
		if err != nil {
			return err.Error()
		}
		res := startActivity(activityAc, client)
		log.Println(res)
		return res
	case constants.ACTION_START_SERVICE:
		return startService(action)
	case constants.APP_CHANGE_STATE:
		return onAppChangeState(client, action)
	case constants.SYSTEM_ACTIVITY_RESULT:
		return onReturnActivityResult(action)
	case constants.APP_TERMINATE:
		return onAppTerminate(client, action)
	case constants.ACTION_APP_RESTART:
		return onAppRestart(client, action)
	case constants.ACTION_APP_CLEAR_DATA:
		return onAppClearData(client, action)
	case constants.ACTION_APP_INSTALLED:
		return initPackage(action)
	case constants.SYSTEM_UPDATE_MODE:
		return activateUpdateMode(client, action)
	case constants.SYSTEM_TERMINATE:
		return terminate(5)
	case constants.SYSTEM_TERMINATE_NOW:
		return terminate(0)
	default:
		return rpc.CreateResponse(rpc.NOT_IMPLEMENTED, "request for action "+action.Id+" was not implemented.")
	}
}

// use to restart the app
func onAppRestart(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	appid := action.PackageId
	if appid == "" {
		return rpc.CreateResponse(rpc.INVALID_CODE, "package id shouldnt be empty")
	}
	app := appman.GetAppConnect(appid, client)
	if app == nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, appid+" app not found")
	}
	if err := app.Restart(); err != nil {
		global.Logger.Error().Err(err).Caller().Msg("failed to restart app " + appid)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed to restart app "+appid)
	}
	return rpc.CreateSuccessResponse("restarted")
}

func onAppClearData(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	appid := action.PackageId
	if appid == "" {
		return rpc.CreateResponse(rpc.INVALID_CODE, "package id shouldnt be empty")
	}
	app := appman.GetAppConnect(appid, client)
	if app == nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, appid+" app not found")
	}
	if err := app.ClearData(); err != nil {
		global.Logger.Error().Err(err).Caller().Msg("failed to clear data for app " + appid)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed to clear data for app "+appid)
	}
	return rpc.CreateSuccessResponse("cleared")
}

func onAppTerminate(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	appid := action.PackageId
	if appid == "" {
		return rpc.CreateResponse(rpc.INVALID_CODE, "package id shouldnt be empty")
	}
	app := appman.GetAppConnect(appid, client)
	if app == nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, appid+" app not found")
	}
	if err := app.Terminate(); err != nil {
		global.Logger.Error().Err(err).Caller().Msg("failed to terminate app " + appid)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "failed to terminate app "+appid)
	}
	return rpc.CreateSuccessResponse("terminated")
}

func onAppChangeState(
	client protocol.ClientInterface,
	action data.Action) interface{} {
	appstate, err := action.DataToAppState()
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, "unable to read app change state value for "+action.PackageId)
	}

	switch {
	case appstate.State == constants.APP_AWAKE || appstate.State == constants.APP_AWAKE_DEBUG:
		var app *appman.AppConnect
		if appstate.State == constants.APP_AWAKE {
			app = appman.GetAppConnect(action.PackageId, client)
		} else {
			app_cwd := appstate.Data.(string)
			app, err = debugging.DebugApp(action.PackageId, app_cwd, client)
			if err != nil {
				logger.GetInstance().Error().Err(err).Caller().Msg("failed debugging app " + action.PackageId)
				return rpc.CreateResponse(rpc.INVALID_CODE, "failed debugging app "+action.PackageId)
			}
		}
		if app != nil {
			return rpc.CreateJsonResponse(200, app.PendingActions)
		} else {
			global.Logger.Warn().Caller().Msg("Trying to awake package " + action.PackageId + " but not installed.")
			return rpc.CreateResponse(rpc.INVALID_CODE, "package not installed")
		}
	case appstate.State == constants.APP_SLEEP:
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
			global.Logger.Error().Caller().Err(err).Msg("Failed to retrieve package with activity " + action.Id)
			return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "StartActivity failed to start "+packageId+" "+err.Error())
		}
		if len(pks) > 0 {
			packageId = pks[0]
		} else {
			global.Logger.Error().Caller().Msg("Cant start activity, package couldn't be found " + action.Id)
			return rpc.CreateResponse(rpc.INVALID_CODE, "cant find package"+packageId+" with action "+action.Id)
		}
	}
	if err := appman.LaunchAppActivity(packageId, client, action); err != nil {
		global.Logger.Error().Caller().Err(err).Msg("Failed to launch activity " + action.Id)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	global.Logger.Debug().Msg("Start activity with " + action.Id + action.DataToString())
	return rpc.CreateSuccessResponse("started")
}

func startService(action data.Action) string {
	pkgid := action.PackageId
	if pkgid == "" {
		return rpc.CreateResponse(rpc.INVALID_CODE, "package id shouldnt be empty")
	}
	err := appman.InitializePackage(pkgid)
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse("started")
}

// when activity returns a result
func onReturnActivityResult(action data.Action) string {
	packageId := action.PackageId
	if packageId == "" {
		return rpc.CreateResponse(rpc.INVALID_CODE, "package id should be the activity whho return the result")
	}
	app := appman.GetAppConnect(packageId, nil)
	if app == nil || app.StartedBy == "" {
		return rpc.CreateResponse(rpc.INVALID_CODE, "cant find the app that would recieve result")
	}
	if !appman.IsAppRunning(app.StartedBy) {
		return rpc.CreateResponse(rpc.INVALID_CODE, "cant find the app that would recieve result")
	}
	// step: start calling the package
	originApp := appman.GetAppConnect(app.StartedBy, nil)
	res, err := originApp.RPC.CallAct(action)
	if err != nil {
		global.Logger.Error().Err(err).Caller().Msg("Failed calling activity result to the calling activity " + app.StartedBy)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse(res)
}

// send RPC to specific package
func sendPackageRPC(pkid string, action data.Action) string {
	if pkid == "" {
		global.Logger.Error().Caller().Msg("No package provided for action " + action.Id)
		return rpc.CreateResponse(rpc.INVALID_CODE, "No package provided for action "+action.Id)
	}
	// step: get package
	app := appman.GetAppConnect(pkid, nil)
	if app == nil {
		global.Logger.Error().Caller().Msg("Unable to send RPC, cant find package " + pkid)
		return rpc.CreateResponse(rpc.INVALID_CODE, "Unable to send RPC, cant find package "+pkid)
	}
	// step: call rpc
	if app.Client == nil {
		global.Logger.Error().Caller().Msg("App is not running " + pkid)
		return rpc.CreateResponse(rpc.INVALID_CODE, "Package "+pkid+" is not loaded yet.")
	}
	res, err := app.RPC.CallAct(action)
	if err != nil {
		global.Logger.Error().Err(err).Caller().Msg("Failed to call RPC for package " + pkid)
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	// step: clean response. remove \"
	if len(res) > 0 {
		res = res[1 : len(res)-1]
	}
	return res
}

// client requested to activate update mode
func activateUpdateMode(client protocol.ClientInterface, action data.Action) string {
	pk := action.DataToString()
	if err := global.Server.EventServer.BroadcastAction(data.NewActionById(constants.BCAST_TERMINATE_N_UPDATE)); err != nil {
		global.Logger.Error().Caller().Err(err).Msg("Failed to broadcast" + constants.BCAST_TERMINATE_N_UPDATE)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	startActivity(data.NewAction(constants.ACTION_APP_SYSTEM_INSTALL, "", pk), nil)
	return rpc.CreateSuccessResponse("Activated")
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
	return rpc.CreateSuccessResponse("Terminated")
}

// initialize and start newly instlled package
func initPackage(action data.Action) string {
	global.Logger.Debug().Msg("Initializing package " + action.PackageId)
	pki := action.PackageId
	if err := appman.InitializePackage(pki); err != nil {
		global.Logger.Error().Err(err).Msg("Failed to initialize package " + pki)
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse("Initialized")
}
