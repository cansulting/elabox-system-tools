package main

import (
	"os"

	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/broadcast"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"

	data2 "github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"

	adata "github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
)

var systemVersion = ""

type MyService struct {
}

func (instance *MyService) OnStart() error {
	if err := broadcast.Init(); err != nil {
		return err
	}

	// register service rpc
	global.AppController.RPC.OnRecieved(global.RETRIEVE_PACKAGES, instance.rpc_retrievePackages)
	global.AppController.RPC.OnRecieved(global.RETRIEVE_PACKAGE, instance.rpc_retrievePackage)
	global.AppController.RPC.OnRecieved(global.INSTALL_PACKAGE, instance.rpc_installPackage)
	global.AppController.RPC.OnRecieved(global.UNINSTALL_PACKAGE, instance.rpc_onuninstall)
	global.AppController.RPC.OnRecieved(global.CANCEL_INSTALL_PACKAGE, instance.rpc_oncancelinstall)
	global.AppController.RPC.OnRecieved(global.RETRIEVE_SYS_VERSION, instance.rpc_onRetrieveSysVersion)
	//go RetrieveAllApps(false)
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

func (instance *MyService) rpc_onRetrieveSysVersion(client protocol.ClientInterface, action data.Action) string {
	// load json file from SYS_INFO_PATH
	if systemVersion == "" {
		contents, err := os.ReadFile(global.SYS_INFO_PATH)
		if err != nil {
			return rpc.CreateResponse(rpc.SYSTEMERR_CODE, "unable to readfile "+err.Error())
		}
		pkg := adata.DefaultPackage()
		if err := pkg.LoadFromBytes(contents); err != nil {
			return rpc.CreateResponse(rpc.INVALID_CODE, err.Error())
		}
		systemVersion = pkg.Version
	}
	return rpc.CreateSuccessResponse(systemVersion)
}

func (instance *MyService) IsRunning() bool {
	return true
}

func (instance *MyService) OnEnd() error {
	return nil
}
