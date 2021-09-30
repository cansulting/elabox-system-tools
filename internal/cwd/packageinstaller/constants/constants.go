package constants

import (
	"ela/foundation/app"
	"ela/foundation/logger"
	"ela/foundation/path"
	"ela/internal/cwd/global"
)

const TERMINATE_TIMEOUT = 5 // seconds to wait for system to terminate
var AppController *app.Controller

// the current logger
var Logger = logger.Init("ela.installer")

// temp path for cache files for
func GetTempPath() string {
	return path.GetCacheDir() + "/custominstall"
}

// temp path for custom installer
func GetCustomInstallerTempPath() string {
	return GetTempPath() + "/" + global.PACKAGEKEY_CUSTOM_INSTALLER
}
