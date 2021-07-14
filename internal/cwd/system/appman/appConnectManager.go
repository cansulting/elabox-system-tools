package appman

import (
	"ela/foundation/event/protocol"
	"log"
)

// currently running processes
var running map[string]*AppConnect = make(map[string]*AppConnect)

func GetAllRunningApps() map[string]*AppConnect {
	return running
}

// use to run process for specific package. return true if success, false if already running
func getAppConnect(packageId string, client protocol.ClientInterface) *AppConnect {
	app, ok := running[packageId]

	// is already running? return false
	if ok {
		if app.client == nil {
			app.client = client
		}
		return app
	}
	pk, _ := RetrievePackage(packageId)
	if pk == nil {
		return nil
	}
	app = newAppConnect(pk)
	app.client = client
	running[packageId] = app
	return app
}

func lookupAppConnect(packageId string) *AppConnect {
	pk, ok := running[packageId]
	// is already running? return false
	if ok {
		return pk
	}
	return nil
}

func RemoveAppConnect(packageId string, terminate bool) {
	app := lookupAppConnect(packageId)
	if app != nil {
		if terminate {
			if err := app.Terminate(); err != nil {
				log.Println("appConnectManager.TerminateAllApp failed terminate "+app.packageId+". Trying force terminate.", err)
				if err := app.ForceTerminate(); err != nil {
					log.Println("appConnectManager.TerminateAllApp failed force terminate ", err)
				}
			}
		}
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
