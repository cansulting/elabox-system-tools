package main

import (
	"ela/foundation/app"
	"ela/foundation/event/data"
)

func InitLibrary() {

}

// subscribe all packages
/*
func subscribeAllApp() {
	appsAction := loadAllAppDefinitions()
	for _, package := range appsAction {
		for _, actionDef := range package.Actions {
			app.GetConnector().Subscribe(
				data.Subscription{Action: actionDef.Action, AppId: package.AppId}, onActionBroadcast
			)}
	}
}*/

func subscribeAppManager() {
	// subscribe to handle launch
	app.GetConnector().Subscribe(
		data.Subscription{Action: ACTION_LAUNCH},
		onActionBroadcast)
}

// callback when a event broadcast an action
func onActionBroadcast(action data.Action) {

}

func loadAllAppDefinitions() []PackageData {
	return nil
}

func Launch() {}
