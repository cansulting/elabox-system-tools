package appman

import (
	"ela/foundation/constants"
	"ela/foundation/event/data"
)

// a client broadcast an action. handle those actions
func onActionResolver(action data.Action) string {
	switch action.Id {
	case constants.ACTION_APP_LAUNCH:
		return launchPackage(action, action.PackageId)
	default:
		// package id was given. then execute the package
		if action.PackageId != "" {
			return launchPackage(action, action.PackageId)
		} else {
			return launcAction(action)
		}
	}
}

func launchPackage(action data.Action, packageId string) string {
	app := getAppConnect(packageId)
	app.pendingActions.AddPendingActivity(action)
	err := app.launch()
	if err != nil {
		return err.Error()
	}
	return ""
}

func launcAction(action data.Action) string {
	pks, err := RetrievePackagesWithActivity(action.Id)
	if err != nil {
		return err.Error()
	}
	if len(pks) > 0 {
		return launchPackage(action, pks[0])
	}

	pks, err = RetrievePackagesWithBroadcast(action.Id)
	if err != nil {
		return err.Error()
	}
	for _, pk := range pks {
		launchPackage(action, pk)
	}
	return ""
}
