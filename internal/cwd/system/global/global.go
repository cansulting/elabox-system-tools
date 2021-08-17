package global

import "ela/server"

var Server *server.Manager

const DB_NAME = "system.dat"

var Running bool = true

const INSTALLER_PKG_ID = "ela.installer"
