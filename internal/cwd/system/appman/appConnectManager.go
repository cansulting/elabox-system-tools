package appman

import (
	"ela/foundation/event/data"
	"ela/foundation/event/protocol"
	registry "ela/registry/app"
	"log"
)

// currently running processes
var running map[string]*AppConnect = make(map[string]*AppConnect)

func GetAllRunningApps() map[string]*AppConnect {
	return running
}

// use to run process for specific package. return true if success, false if already running
func GetAppConnect(packageId string, client protocol.ClientInterface) *AppConnect {
	app, ok := running[packageId]

	// is already running? return false
	if ok {
		if client != nil {
			app.Client = client
		}
		return app
	}
	// retrieve if already exist
	pk, _ := registry.RetrievePackage(packageId)
	if pk == nil {
		return nil
	}
	// create service if exist
	//var service *ServiceConnect = nil
	//if pk.HasServices() {
	//service = onServiceOpen(client, pk.PackageId)
	//}
	app = newAppConnect(pk, client)
	running[packageId] = app
	return app
}

func LookupAppConnect(packageId string) *AppConnect {
	pk, ok := running[packageId]
	// is already running? return false
	if ok {
		return pk
	}
	return nil
}

// use to check if app is currently running or not
func IsAppRunning(packageId string) bool {
	app := LookupAppConnect(packageId)
	if app == nil {
		return false
	}
	return true
}

func RemoveAppConnect(packageId string, terminate bool) {
	app := LookupAppConnect(packageId)
	if app != nil {
		if terminate {
			if err := app.Terminate(); err != nil {
				log.Println("appConnectManager.TerminateAllApp failed terminate "+app.PackageId+". Trying force terminate.", err)
				if err := app.ForceTerminate(); err != nil {
					log.Println("appConnectManager.TerminateAllApp failed force terminate ", err)
				}
			}
		}
		// close service
		// if app.Service != nil {
		// 	onServiceClose(packageId)
		// }
		delete(running, packageId)
	}
}

func TerminateAllApp() {
	log.Println("appConnectManager.TerminateAllApp() started")
	running := GetAllRunningApps()
	for pkid := range running {
		RemoveAppConnect(pkid, true)
	}
}

// this launches the activity
func LaunchAppActivity(
	packageId string,
	caller protocol.ClientInterface,
	pendingActivity data.Action) error {
	// start launching the activity
	app := GetAppConnect(packageId, nil)
	app.PendingActions.AddPendingActivity(pendingActivity)
	err := app.Launch()
	if err != nil {
		return err
	}
	return nil
}
