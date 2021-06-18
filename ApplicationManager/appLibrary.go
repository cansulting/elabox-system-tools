package main

import (
	base "ela.services/Base"
	client "ela.services/ServiceCenterLib"
)

func InitLibrary() {

}

// subscribe all packages
func subscribeAllApp() {
	appsAction := loadAllAppDefinitions()
	for _, app := range appsAction {
		for _, actionDef := range app.Actions {
			base.GetConnector().Subscribe(
				client.SubscriptionData{Action: actionDef.Action, AppId: app.AppId}, onActionBroadcast)
		}
	}
}

func subscribeAppManager() {
	// subscribe to handle launch
	base.GetConnector().Subscribe(
		client.SubscriptionData{Action: ACTION_LAUNCH},
		onActionBroadcast)
}

// callback when a client broadcast an action
func onActionBroadcast(action client.ActionData) {

}

func loadAllAppDefinitions() []PackageData {
	return nil
}

func Launch() {}
