package global

import "ela/foundation/event/protocol"

var Connector protocol.ConnectorServer

const DB_NAME = "system.dat"

var Running bool = true

const INSTALLER_PKG_ID = "ela.installer"
