// this file handles broadcast for normal installation

package broadcast

import (
	"strconv"

	appc "github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/packageinstaller/constants"
)

// broadcast progress
// @param pkgId - the package id currently installing
// @param progress - the progress of the package 0 -100
func UpdateProgress(pkgId string, progress int) {
	value := `{"packageId":"` + pkgId + `","progress":` + strconv.Itoa(progress) + `}`
	_, err := constants.AppController.RPC.CallBroadcast(data.NewAction(
		constants.INSTALLER_PROGRESS,
		constants.PKG_ID,
		value))
	if err != nil {
		constants.Logger.Error().Err(err).Msg("failed to broadcast progress + " + strconv.Itoa(int(progress)))
	}
}

// broadcast state changed
// @param pkgId - the package id currently installing
// @param state - the state of the installation
func UpdateSystem(pkgId string, status InstallState) {
	value := `{"packageId":"` + pkgId + `","status":"` + string(status) + `"}`
	_, err := constants.AppController.RPC.CallBroadcast(data.NewAction(
		constants.INSTALLER_STATE_CHANGED,
		constants.PKG_ID,
		value))
	if err != nil {
		constants.Logger.Error().Err(err).Msg("failed to broadcast status update")
	}
}

// notify system that the installation is complete for specific package
func OnPackageInstalled(pki string) error {
	_, err := constants.AppController.RPC.CallSystem(data.NewAction(
		appc.ACTION_APP_INSTALLED, pki, nil,
	))
	if err != nil {
		return err
	}
	return nil
}

// broadcast error
// @param pkgId - the package id currently installing
func Error(pkgId string, code int, val string) {
	value := `{"packageId":"` + pkgId + `","code":` + strconv.Itoa(code) + `,"error":"` + val + `"}`
	_, err := constants.AppController.RPC.CallBroadcast(data.NewAction(
		constants.INSTALLER_ERROR,
		constants.PKG_ID,
		value))
	if err != nil {
		constants.Logger.Error().Err(err).Msg("failed to broadcast error")
	}
}
