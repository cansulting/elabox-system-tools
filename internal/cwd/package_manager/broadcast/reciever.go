// callback when recieved broadcast actions

package broadcast

import (
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/event/protocol"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
)

// register broadcast recievers
func registerRecievers() error {
	// register for
	if err := global.RPC.OnRecievedFromPackage(
		global.InstallerId,
		global.INSTALLER_PROGRESS,
		onRecievedInstallerProgress); err != nil {
		return err
	}
	if err := global.RPC.OnRecievedFromPackage(
		global.InstallerId,
		global.INSTALLER_STATE_CHANGE,
		onRecievedInstallerStateChanged); err != nil {
		return err
	}
	if err := global.RPC.OnRecievedFromPackage(
		global.InstallerId,
		global.INSTALLER_ERROR,
		onRecievedInstallerError); err != nil {
		return err
	}
	return nil
}

// callback from installer when it's progress changed
func onRecievedInstallerProgress(client protocol.ClientInterface, action data.Action) string {
	// step: parse data
	dataAc, err := action.DataToMap()
	if err != nil {
		logger.GetInstance().Error().Caller().Err(err).Msg("failed to parse action data")
		return ""
	}
	currentPackage, ok := dataAc["packageId"].(string)
	if !ok || currentPackage == "" {
		logger.GetInstance().Error().Caller().Msg("packageId is not string")
		return ""
	}
	progress, ok := dataAc["progress"].(float64)
	if !ok {
		logger.GetInstance().Error().Caller().Msg("failed to parse progress")
		return ""
	}
	OnInstallerProgress(currentPackage, int(progress))
	return ""
}

// callback when state changed on installer
func onRecievedInstallerStateChanged(client protocol.ClientInterface, action data.Action) string {
	dataAc, err := action.DataToMap()
	if err != nil {
		logger.GetInstance().Error().Caller().Err(err).Msg("failed to parse action data")
		return ""
	}
	currentPk := dataAc["packageId"].(string)
	status := dataAc["status"].(string)
	OnInstallerStateChanged(currentPk, PkInstallerState(status))
	return ""
}

// callback when installer got an issue
func onRecievedInstallerError(client protocol.ClientInterface, action data.Action) string {
	dataAc, err := action.DataToMap()
	if err != nil {
		logger.GetInstance().Error().Caller().Err(err).Msg("failed to parse action data")
		return ""
	}
	currentPk := dataAc["packageId"].(string)
	code := dataAc["code"].(float64)
	errmsg := dataAc["error"].(string)
	OnInstallerError(currentPk, int(code), errmsg)
	return ""
}
