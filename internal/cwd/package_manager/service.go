package main

import (
	"encoding/json"

	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/broadcast"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/system_updater"

	data2 "github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
)

type MyService struct {
}

func (instance *MyService) OnStart() error {
	if err := broadcast.Init(); err != nil {
		return err
	}
	if err := system_updater.Init(); err != nil {
		return err
	}

	// register service rpc
	global.AppController.RPC.OnRecieved(global.RETRIEVE_PACKAGES, instance.rpc_retrievePackages)
	global.AppController.RPC.OnRecieved(global.RETRIEVE_PACKAGE, instance.rpc_retrievePackage)
	global.AppController.RPC.OnRecieved(global.INSTALL_PACKAGE, instance.rpc_installPackage)
	global.AppController.RPC.OnRecieved(global.UNINSTALL_PACKAGE, instance.rpc_onuninstall)
	global.AppController.RPC.OnRecieved(global.CANCEL_INSTALL_PACKAGE, instance.rpc_oncancelinstall)
	global.AppController.RPC.OnRecieved(global.RETRIEVE_SYS_VERSION, instance.rpc_onRetrieveSysVersion)

	// system updater
	global.AppController.RPC.OnRecieved(global.AC_DOWNLOAD_UPDATE, instance.rpc_onDownloadUpdate)
	global.AppController.RPC.OnRecieved(global.AC_INSTALL_UPDATE, instance.rpc_onInstallUpdate)
	return nil
}

// callback RPC
func (instance *MyService) rpc_retrievePackages(client protocol.ClientInterface, action data.Action) string {
	dmap, _ := action.DataToMap()
	includeHidden := true
	if dmap["include_hidden"] != nil && dmap["include_hidden"].(bool) == false {
		includeHidden = false
	}
	apps, err := RetrieveAllApps(includeHidden)
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	return rpc.CreateJsonResponse(rpc.SUCCESS_CODE, apps)
}

func (instance *MyService) rpc_retrievePackage(client protocol.ClientInterface, action data.Action) string {
	app, err := RetrieveApp(action.PackageId, "")
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	return rpc.CreateJsonResponse(rpc.SUCCESS_CODE, app)
}

func (instance *MyService) rpc_installPackage(client protocol.ClientInterface, action data.Action) string {
	// parse rpc params
	dataM, err := action.DataToMap()
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	definition := data2.InstallDef{}.FromMap(dataM["definition"].(map[string]interface{}))
	var dependencies []data2.InstallDef = nil
	if dataM["dependencies"] != nil {
		dependencies = dataM["dependencies"].([]data2.InstallDef)
	}
	// download and install
	err = DownloadInstallApp(definition, dependencies)
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse("started")
}

func (instance *MyService) rpc_onuninstall(client protocol.ClientInterface, action data.Action) string {
	err := UninstallApp(action.PackageId)
	if err != nil {
		return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse("started")
}

func (instance *MyService) rpc_oncancelinstall(client protocol.ClientInterface, action data.Action) string {
	CancelInstall(action.PackageId)
	return rpc.CreateSuccessResponse("cancelled")
}

var tmpMap = make(map[string]interface{})

func (instance *MyService) rpc_onRetrieveSysVersion(client protocol.ClientInterface, action data.Action) string {
	curver := system_updater.GetCurrentSysVersion()
	latver := system_updater.GetLatestSysVersion()
	tmpMap["current"] = curver
	tmpMap["latest"] = latver
	content, err := json.Marshal(tmpMap)
	if err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse(string(content))
}

func (instance *MyService) rpc_onRetrieveUpdateData(client protocol.ClientInterface, action data.Action) string {
	return ""
}

func (*MyService) rpc_onDownloadUpdate(client protocol.ClientInterface, action data.Action) string {
	if err := system_updater.DownloadLatest(); err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse("downloading")
}

func (*MyService) rpc_onInstallUpdate(client protocol.ClientInterface, action data.Action) string {
	if err := system_updater.InstallUpdate(); err != nil {
		return rpc.CreateResponse(rpc.SYSTEMERR_CODE, err.Error())
	}
	return rpc.CreateSuccessResponse("installing")
}

func (instance *MyService) IsRunning() bool {
	return true
}

func (instance *MyService) OnEnd() error {
	return nil
}
