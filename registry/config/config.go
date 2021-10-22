package config

import (
	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

const DB_NAME = "system.dat"

var DB_DIR = path.GetSystemAppDirData(constants.SYSTEM_SERVICE_ID)
