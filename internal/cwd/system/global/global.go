// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package global

import (
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/server"
)

var Server *server.Manager               // the server manager that handles event and web
var Running = true                       // true if this system is currently running
const INSTALLER_PKG_ID = "ela.installer" // package id of installer
const SYSTEM_PKID = "ela.system"         //
const RUN_STARTUPAPPS = true             // true if system runs startup apps
var Logger = logger.Init(SYSTEM_PKID)

const APP_TERMINATE_COUNTDOWN = 3 // number of seconds to wait before terminating an app
const CONFIG_ENV = "config"       // env name for config, value is "1" if elabox was already configured
const DEFAULT_DASHBOARD = "ela.dashboard"

