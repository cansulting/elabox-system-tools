package global

import (
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/server"
)

var Server *server.Manager               // the server manager that handles event and web
var Running bool = true                  // true if this system is currently running
const INSTALLER_PKG_ID = "ela.installer" // package id of installer
const SYSTEM_PKID = "ela.system"
const RUN_STARTUPAPPS = true // true if system runs startup apps
var Logger = logger.Init(SYSTEM_PKID)
