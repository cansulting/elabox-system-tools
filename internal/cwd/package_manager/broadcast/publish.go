package broadcast

import (
	"dashboard/package_manager/data"
	"dashboard/package_manager/global"
	"log"
	"strconv"

	sdata "github.com/cansulting/elabox-system-tools/foundation/event/data"
)

func PublishInstallProgress(progress uint, itemId string) error {
	val := `{"progress":` + strconv.Itoa(int(progress)) + `,"packageId":"` + itemId + `"}`
	_, err := global.RPC.CallBroadcast(sdata.NewAction(global.INSTALL_PROGRESS, global.PackageId, val))
	return err
}

func PublishNewUpdateAvailable(updates []*data.PackageListingCache) error {
	log.Println("updates found: " + strconv.Itoa((len(updates))))
	// create json from updates
	var jsonStr string = "["
	for _, v := range updates {
		if len(jsonStr) > 1 {
			jsonStr += ","
		}
		jsonStr += `{"id":"` + v.Id + `","name":"` + v.Name + `","build":` + strconv.Itoa(v.Build) + `"}`
	}
	jsonStr += "]"
	_, err := global.RPC.CallBroadcast(sdata.NewAction(global.UPDATE_AVAILABLE, global.PackageId, jsonStr))
	return err
}

func PublishError(pkid string, code int, msg string) error {
	val := `{"code":` + strconv.Itoa(code) + `,"error":"` + msg + `", "packageId":"` + pkid + `"}`
	_, err := global.RPC.CallBroadcast(sdata.NewAction(global.BROADCAST_ERROR, global.PackageId, val))
	return err
}

// publishes change state event for specific package
// pkid: which package that was changed
func PublishInstallState(pkid string, state global.AppStatus) error {
	val := `{"packageId":"` + pkid + `","status":"` + string(state) + `"}`
	_, err := global.RPC.CallBroadcast(sdata.NewAction(global.INSTALL_STATE, global.PackageId, val))
	return err
}
