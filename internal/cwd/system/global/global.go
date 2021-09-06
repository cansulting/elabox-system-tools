package global

import "ela/server"

var Server *server.Manager               // the server manager that handles event and web
var Running bool = true                  // true if this system is currently running
const INSTALLER_PKG_ID = "ela.installer" // package id of installer
const RUN_STARTUPAPPS = true             // true if system runs startup apps
