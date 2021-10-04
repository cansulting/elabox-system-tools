package constants

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/internal/cwd/global"
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
